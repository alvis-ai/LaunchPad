package main

import (
	"database/sql"
	"flag"
	"fmt"
	"launchpad/config"
	"os"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	// 解析命令行参数
	sqlQuery := flag.String("sql", "", "SQL查询语句")
	dbPassword := flag.String("password", "", "数据库密码（可选，覆盖配置文件中的密码）")
	flag.Parse()

	// 如果没有提供SQL，使用默认查询
	query := *sqlQuery
	if query == "" {
		// 默认查询：查看所有表
		query = "SHOW TABLES"
		fmt.Println("未提供SQL查询，使用默认查询: SHOW TABLES")
		fmt.Println("使用方法: go run temp/query.go -sql 'SELECT * FROM product_contract'")
		fmt.Println("---")
	}

	// 加载配置
	config.LoadConfig()

	// 确定使用的密码（优先使用命令行参数，其次使用环境变量，最后使用配置文件）
	password := config.AppConfig.Database.Password
	if *dbPassword != "" {
		password = *dbPassword
		fmt.Println("使用命令行指定的密码")
	} else if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
		password = envPassword
		fmt.Println("使用环境变量 DB_PASSWORD 指定的密码")
	}

	// 显示数据库配置信息（隐藏密码）
	fmt.Printf("数据库配置: %s@%s:%s/%s\n",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name)

	// 直接创建数据库连接（不依赖 database.InitializeDB）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.User,
		password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("打开数据库连接失败: %v", err)
	}
	defer sqlDB.Close()

	// 测试数据库连接是否有效
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v\n\n请检查:\n1. 数据库服务是否运行 (docker ps | grep mysql)\n2. 配置文件中的数据库连接信息是否正确\n3. 数据库用户权限是否正确\n4. 密码是否正确\n\n提示: 可以使用 -password 参数覆盖密码，例如:\n   go run temp/query.go -password 123456\n   或设置环境变量: export DB_PASSWORD=123456", err)
	}

	fmt.Println("✓ 数据库连接成功")

	// 执行查询
	fmt.Printf("\n执行 SQL: %s\n", query)
	rows, err := sqlDB.Query(query)
	if err != nil {
		log.Fatalf("执行SQL查询失败: %v", err)
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("获取列名失败: %v", err)
	}

	// 读取所有数据
	var allRows [][]string
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	rowCount := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("扫描行数据失败: %v", err)
		}

		// 将每行数据转换为字符串
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				// 将 []byte 转换为字符串
				if b, ok := val.([]byte); ok {
					row[i] = string(b)
				} else {
					row[i] = fmt.Sprintf("%v", val)
				}
			}
		}
		allRows = append(allRows, row)
		rowCount++
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("遍历行时出错: %v", err)
	}

	// 计算每列的最大宽度
	colWidths := make([]int, len(columns))
	for i, col := range columns {
		colWidths[i] = len(col)
	}
	for _, row := range allRows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 限制最大列宽，避免过宽
	maxWidth := 50
	for i := range colWidths {
		if colWidths[i] > maxWidth {
			colWidths[i] = maxWidth
		}
	}

	// 打印表头
	fmt.Println("\n查询结果:")
	printSeparator(colWidths)
	printRow(columns, colWidths)
	printSeparator(colWidths)

	// 打印数据行
	for _, row := range allRows {
		printRow(row, colWidths)
	}

	if rowCount > 0 {
		printSeparator(colWidths)
	}
	fmt.Printf("共查询到 %d 行数据\n", rowCount)
}

// printRow 打印表格的一行
func printRow(row []string, widths []int) {
	for i, cell := range row {
		// 如果内容超过最大宽度，截断并添加省略号
		display := cell
		if len(display) > widths[i] {
			display = display[:widths[i]-3] + "..."
		}
		fmt.Printf("| %-*s ", widths[i], display)
	}
	fmt.Println("|")
}

// printSeparator 打印分隔线
func printSeparator(widths []int) {
	for _, width := range widths {
		fmt.Print("+")
		for i := 0; i < width+2; i++ {
			fmt.Print("-")
		}
	}
	fmt.Println("+")
}
