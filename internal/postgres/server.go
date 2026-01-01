package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"al.essio.dev/pkg/shellescape"
	"github.com/kellegous/poop"
	_ "github.com/lib/pq"
)

const (
	defaultPgPath = "/usr/lib/postgresql/14" // TODO(knorton): make this configurable
	pgUser        = "postgres"
)

type Server struct {
	dataDir   string
	pgPath    string
	pgVersion int
}

func (s *Server) pgDataDir() string {
	return filepath.Join(s.dataDir, strconv.Itoa(s.pgVersion))
}

func Start(
	ctx context.Context,
	dataDir string,
	opts ...Option,
) (*Server, error) {
	version, err := getPgVersion(ctx)
	if err != nil {
		return nil, poop.Chain(err)
	}

	s := &Server{
		dataDir:   dataDir,
		pgPath:    defaultPgPath,
		pgVersion: version,
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.ensureVersionDir(); err != nil {
		return nil, poop.Chain(err)
	}

	if err := s.initDB(ctx); err != nil {
		return nil, poop.Chain(err)
	}

	if err := s.start(ctx); err != nil {
		return nil, poop.Chain(err)
	}

	return s, nil
}

func (s *Server) EnsureDatabase(
	ctx context.Context,
	name string,
	username string,
	password string,
) error {
	// TODO(knorton): limit name, username, password to valid characters
	q := fmt.Sprintf(`
		CREATE DATABASE %s;
		CREATE USER %s WITH ENCRYPTED PASSWORD '%s';
		GRANT ALL PRIVILEGES ON DATABASE %s TO %s;`,
		name, username, password, name, username)
	if err := psql(ctx, pgUser, q); err != nil {
		return err
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable",
			name,
			username,
			password))
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return suCommand(
		ctx,
		pgUser,
		filepath.Join(s.pgPath, "bin/pg_ctl"),
		"-D", s.dataDir,
		"stop").Run()
}

func (s *Server) ensureVersionDir() error {
	// if the version file isn't there, we don't need to do anything.
	if _, err := os.Stat(filepath.Join(s.dataDir, "PG_VERSION")); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	// now, we need to move all files/dirs from the root data dir to the version dir.
	dataDir := s.pgDataDir()

	// ensure the version dir exists.
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			return poop.Chain(err)
		}

		uid, gid, err := getUser(pgUser)
		if err != nil {
			return poop.Chain(err)
		}

		if err := os.Chown(dataDir, uid, gid); err != nil {
			return poop.Chain(err)
		}
	}

	files, err := os.ReadDir(s.dataDir)
	if err != nil {
		return poop.Chain(err)
	}

	versionName := strconv.Itoa(s.pgVersion)
	for _, file := range files {

		if file.Name() == versionName {
			continue
		}

		if err := os.Rename(
			filepath.Join(s.dataDir, file.Name()),
			filepath.Join(dataDir, file.Name()),
		); err != nil {
			return poop.Chain(err)
		}
	}

	return nil
}

func (s *Server) initDB(ctx context.Context) error {
	dataDir := s.pgDataDir()
	versionPath := filepath.Join(dataDir, "PG_VERSION")
	if _, err := os.Stat(versionPath); err == nil {
		return nil
	}

	uid, gid, err := getUser(pgUser)
	if err != nil {
		return poop.Chain(err)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return poop.Chain(err)
	}

	if err := os.Chown(dataDir, uid, gid); err != nil {
		return poop.Chain(err)
	}

	c := suCommand(
		ctx,
		pgUser,
		filepath.Join(s.pgPath, "bin/initdb"),
		"-D",
		dataDir)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return poop.Chain(c.Run())
}

func (s *Server) waitForReady(
	ctx context.Context,
) error {
	for {
		c := suCommand(
			ctx,
			pgUser,
			filepath.Join(s.pgPath, "bin/pg_isready"),
			"-q")
		if err := c.Run(); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func (s *Server) start(ctx context.Context) error {
	c := suCommand(
		ctx,
		pgUser,
		filepath.Join(s.pgPath, "bin/pg_ctl"),
		"-D",
		s.pgDataDir(),
		"start")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		return err
	}

	if err := s.waitForReady(ctx); err != nil {
		return err
	}

	return nil
}

func suCommand(
	ctx context.Context,
	user string,
	args ...string,
) *exec.Cmd {
	return exec.CommandContext(
		ctx,
		"su",
		"-",
		user,
		"-c",
		shellescape.QuoteCommand(args),
	)
}

func getPgVersion(ctx context.Context) (int, error) {
	c := exec.CommandContext(ctx, "pg_config", "--version")
	var buf bytes.Buffer
	c.Stdout = &buf

	if err := c.Run(); err != nil {
		return 0, poop.Chain(err)
	}

	version := strings.TrimPrefix(buf.String(), "PostgreSQL ")
	major, _, ok := strings.Cut(version, ".")
	if !ok {
		return 0, poop.Newf("invalid version: %s", version)
	}

	return strconv.Atoi(major)
}

func psql(
	ctx context.Context,
	user string,
	sql string,
) error {
	c := suCommand(
		ctx,
		user,
		filepath.Join(defaultPgPath, "bin/psql"))
	c.Stdin = bytes.NewBufferString(sql)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func getUser(username string) (int, int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}
