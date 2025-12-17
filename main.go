package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	c "Emissary/configs"
	l "Emissary/ldap"
	t "Emissary/tree"
)

func init() {
	// Инициализирую ключи программы
	flag.StringVar(&c.Login, "l", "", c.CLogin)
	flag.StringVar(&c.Pass, "p", "", c.CPass)
	flag.StringVar(&c.LdapServer, "s", "", c.CLdapServer)
	flag.StringVar(&c.LdapBaseDn, "b", "", c.CLdapBaseDn)
	flag.StringVar(&c.LdapFilter, "f", "", c.CLdapFilter)
	flag.StringVar(&c.LdapSkipRegexp, "r", "", c.CLdapSkipRegexp)
	flag.StringVar(&c.Output, "o", "Emissary.html", c.COutput)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), c.CUsage+"\n\n")
		flag.PrintDefaults()
	}
}

func main() {
	// Запрашиваю аргументы командной строки и проверяю их наличие
	flag.Parse()
	if c.Login == "" || c.Pass == "" ||
		c.LdapServer == "" || c.LdapBaseDn == "" || c.LdapFilter == "" {
		flag.Usage()
		panic(c.CErr)
	}

	// Запрашиваю пользователей из LDAP
	users, err := l.GetUsersFromLdap()
	if err != nil {
		slog.Error("Ошибка при запросе в LDAP", slog.Any("err", err))
		os.Exit(1)
	}

	// Формирую дерево пользователей по отделам
	tree := t.MakeTree(users)

	// Вывожу дерево пользователей в файл
	if err := t.PrintTree(&tree); err != nil {
		slog.Error("Ошибка при выводе справочника в файл", slog.Any("err", err))
		os.Exit(1)
	}
}
