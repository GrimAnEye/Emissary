package ldap

import (
	"log/slog"
	"regexp"
	"strings"

	c "Emissary/configs"

	ldap "github.com/go-ldap/ldap/v3"
)

// Список возвращаемых атрибутов
var requestAttributes = []string{
	"description",
	"telephoneNumber",
	"otherTelephone",
	"mobile",
	"otherMobile",
	"mail",
}

// GetUsersFromLdap - запрашивает пользователей из LDAP и возвращает их список
func GetUsersFromLdap() (users []c.LdapAttributes, err error) {
	slog.Debug("Подключаюсь к LDAP",
		slog.String("ldapServer", c.LdapServer),
		slog.String("ldapUser", c.User),
		slog.String("ldapFilter", c.LdapFilter),
		slog.String("ldapBase", c.LdapBaseDn),
		slog.Any("ldapRequestAttributes", requestAttributes),
		slog.String("ldapSkipRegexp", c.LdapSkipRegexp),
	)

	// Формирование подключения к ldap
	l, err := ldap.DialURL(c.LdapServer)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	// Использование УЗ для подключения
	err = l.Bind(c.User, c.Pass)
	if err != nil {
		return nil, err
	}

	// Создание запроса в AD
	searchRequest := ldap.NewSearchRequest(
		c.LdapBaseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		c.LdapFilter,
		requestAttributes,
		nil,
	)
	slog.Debug("Запрашиваю пользователей")

	// Выполнение запроса в AD
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	slog.Info("Успешный запрос ldap", slog.Int("objectsCount", len(sr.Entries)))

	// Инициализирую массив пользователей и складываю в него данные
	users = make([]c.LdapAttributes, 0)

	// Подготовка регурярки для дополнительного пропуска пользователей
	re := &regexp.Regexp{}
	if c.LdapSkipRegexp != "" {
		re, err = regexp.Compile(c.LdapSkipRegexp)
		if err != nil {
			return nil, err
		}
	}
	// Прохожусь по всем сущностям и формирую список пользователей
	for _, entry := range sr.Entries {
		if c.LdapSkipRegexp != "" {
			skips := re.FindAllStringSubmatch(entry.DN, -1)
			if len(skips) > 0 {
				slog.Info("Пропускаю", slog.Any("regStrings", skips), slog.String("user", entry.DN))
				continue
			}
		}

		var user c.LdapAttributes

		// Должность пользователя
		user.Position = entry.GetAttributeValue("description")

		// Суммирую все контактные данные пользователя
		for _, attr := range requestAttributes {
			user.Contacts = append(user.Contacts, entry.GetAttributeValues(attr)...)
		}

		// Разбираю DN строку для получения имени и каталогов пользователя
		dn := strings.TrimSuffix(entry.DN, ","+c.LdapBaseDn) // Отрезаю BaseDn
		for strings.Contains(dn, "OU=") {

			indexOU := strings.LastIndex(dn, ",OU=")    // Расположение последнего каталога в строке
			user.OUs = append(user.OUs, dn[indexOU+4:]) // Прибавляю каталог к массиву каталогов пользователя
			dn = dn[:indexOU]                           // Отрезаю данный каталог от строки
		}

		user.Name = strings.TrimPrefix(dn, "CN=") // Получаю имя пользователя

		users = append(users, user)
	}
	return users, nil
}
