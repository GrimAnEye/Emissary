package configs

type (
	// ldapAttributes - атрибуты пользователя в AD
	LdapAttributes struct {
		Name     string   // Имя пользователя
		Position string   // Должность пользователя
		OUs      []string // Каталоги расположения пользователя
		Contacts []string // Контактные данные
	}

	// branchAttributes - атрибуты пользователя для вывода
	BranchAttributes struct {
		Position string // Должность пользователя
		Contacts string // Контактные данные пользователя
	}

	// branch - структура для составления дерева пользователей, в зависимости от их расположения по отделам
	Branch struct {
		Users  map[string]BranchAttributes // Пользователи на ветке
		Branch map[string]*Branch          // Подветви текущей ветки
	}
)

var (
	// Запрашиваемые переменные от пользователя

	User           string // Имя пользователя для доступа к LDAP
	Pass           string // Пароль пользователя для доступа к LDAP
	LdapServer     string // Адрес LDAP сервера
	LdapBaseDn     string // Базовый каталог LDAP для поиска
	LdapFilter     string // LDAP фильтр для поиска пользователей
	LdapSkipRegexp string // Пропускает пользователя, если его DN удовлетворяет регулярному выражению
	Output         string // Путь до файла вывода телефонного справочника
	Log            string // Уровень журнала - error,[info],debug

)

const (
	// Описание запрашиваемых переменных

	CUser           string = "Имя пользователя для доступа к LDAP: user@domain.com"
	CPass           string = "Пароль пользователя для доступа к LDAP" //gosec:disable G101 -- Ошибка анализатора
	CLdapServer     string = "Адрес LDAP сервера в формате: ldaps://ldap.example.com:389"
	CLdapBaseDn     string = "Базовый каталог LDAP для поиска: OU=Users,DC=exa,DC=example,DC=com"
	CLdapFilter     string = "LDAP фильтр для поиска пользователей: (objectClass=organizationalPerson)"
	CLdapSkipRegexp string = "Пропускает пользователя, если его DN удовлетворяет регулярному выражению"
	COutput         string = "Путь до файла вывода телефонного справочника"
	CLog            string = "Уровень журнала - error,[info],debug"

	CUsage string = `Emissary - программа офисной техподдержки,
	формирующая HTML страницу "телефонного" справочника пользователей
	с их контактными данными, для offline использования`
	CErr string = `Требуется указать данные для всех ключей программы! Справка "-h"`
)
