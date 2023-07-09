package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alessio/shellescape"
	_ "github.com/lib/pq"
)

const (
	defaultPgPath = "/usr/lib/postgresql/14"
	pgUser        = "postgres"
)

type Server struct {
	dataDir string
	pgPath  string
}

func Start(
	ctx context.Context,
	dataDir string,
	opts ...Option,
) (*Server, error) {
	s := &Server{
		dataDir: dataDir,
		pgPath:  defaultPgPath,
	}

	for _, opt := range opts {
		opt(s)
	}

	if err := s.initDB(ctx); err != nil {
		return nil, err
	}

	if err := s.start(ctx); err != nil {
		return nil, err
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

func (s *Server) initDB(ctx context.Context) error {
	versionPath := filepath.Join(s.dataDir, "PG_VERSION")
	if _, err := os.Stat(versionPath); err == nil {
		return nil
	}

	uid, gid, err := getUser(pgUser)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return err
	}

	if err := os.Chown(s.dataDir, uid, gid); err != nil {
		return err
	}

	c := suCommand(
		ctx,
		pgUser,
		filepath.Join(s.pgPath, "bin/initdb"),
		"-D",
		s.dataDir)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	return c.Run()
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
		s.dataDir,
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
