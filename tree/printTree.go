package tree

import (
	"fmt"
	"os"

	c "Emissary/configs"
)

// PrintTree - печатает дерево пользователей в файл
func PrintTree(tree *c.Branch) error {
	file, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := initFile(file, true); err != nil {
		return err
	}
	if err := printBranch(tree, file, true); err != nil {
		return err
	}
	if err := initFile(file, false); err != nil {
		return err
	}

	fmt.Println("Вывод данных в файл окончен. Имя справочника: ", c.Output)

	return nil
}

// initFile - запись в файл вывода стартовых и конечных значений
func initFile(file *os.File, isStart bool) error {
	// Шаблоны для конца и начала выходного документа
	top := `<html>
	<head>
		<meta charset="utf-8">
				<title>Телефонный справочник</title>
				<style>
				
				table {
					width: 100%;
					border-collapse: collapse;
				}
				
				td {
					text-align: center;
					padding: 0px;
				}
				
				.governance {
					padding: 2px 4px 2px 4px;
					background-color: rgb(203, 220, 237);
					font-size: 18px;
				}
				
				.groups {
					padding: 3px 3px 3px 3px;
					border-bottom: 1px solid grey;
					background-color: rgb(211, 211, 211);
				}
				
				.member {
					border-bottom: 1px solid;
				}
				
				</style>
				</head>
				<body><table>` + "\n"

	bottom := `</table></body></html>`

	if _, err := file.WriteString(
		// В зависимости от того, начало это документа или конец - записывать разные шаблоны
		func() string {
			if isStart {
				return top
			} else {
				return bottom
			}
		}()); err != nil {
		return err
	}
	return nil
}

// printBranch - перебирает подвертви, печатает их название и вызывает печать пользователей
func printBranch(branch *c.Branch, file *os.File, isStart bool) error {
	if branch.Users != nil {
		// Если есть пользователи, создаю блок заголовка
		_, err := file.WriteString(
			"<tr>" +
				"<td class=\"fio\">ФИО</td>" +
				"<td class=\"contact\">Контакты</td>" +
				"<td class=\"position\">Должность</td>" +
				"</tr>\n")
		if err != nil {
			return err
		}
		// Начинаю выводить в него пользователей
		printUsers(&branch.Users, file)
	}

	if branch.Branch != nil {
		// Если есть подветви, создаю блок названия отдела
		for subBranch := range branch.Branch {

			if _, err := file.WriteString(

				fmt.Sprintf("<thead class=\"%s\"><tr><th colspan=\"3\">%s</th></tr></thead>\n",

					// В зависимости, главный это каталог или дочерний, назначать разные стили
					func() string {
						if isStart {
							return "governance"
						} else {
							return "groups"
						}
					}(), subBranch,
				),
			); err != nil {
				return err
			}
			// Ухожу на глубину
			if err := printBranch(branch.Branch[subBranch], file, false); err != nil {
				return err
			}

		}
	}

	return nil
}

// printUsers - печатает заголовок-аннотацию и пользователя
func printUsers(users *map[string]c.BranchAttributes, file *os.File) error {
	var out string = ""

	out += "<tbody>"
	// Перебираю пользователей в карте
	for userName, attributes := range *users {
		out += fmt.Sprintf(
			"<tr class=\"member\">"+
				"<td>%s</td>"+
				"<td>%s</td>"+
				"<td>%s</td>"+
				"</tr>\n",
			userName,
			attributes.Contacts,
			attributes.Position)
	}
	out += "</tbody>\n"

	// Записываю строку с данными пользователя
	if _, err := file.WriteString(out); err != nil {
		return err
	}
	return nil
}
