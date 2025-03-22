package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Структура для відображення результатів у шаблоні
type InputData struct {
	PowerMW float64 // Потужність у МВт
	Cost    float64 // Ціна, грн/кВт·год
	Portion float64 // Частка енергії, що не обкладається штрафом (0..1)

	Revenue float64 // Доходи (тис. грн)
	Penalty float64 // Штраф (тис. грн)
	Profit  float64 // Прибуток (тис. грн)
}

func main() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

// Головна сторінка
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Зчитування даних із форми
		r.ParseForm()
		powerMW, _ := strconv.ParseFloat(r.FormValue("powerMW"), 64) // МВт
		cost, _ := strconv.ParseFloat(r.FormValue("cost"), 64)       // грн/кВт·год
		portion, _ := strconv.ParseFloat(r.FormValue("portion"), 64) // частка 0..1

		// Загальний обсяг енергії за добу (МВт·год)
		totalEnergy := powerMW * 24.0

		// Доходи (грн) з частини, яка не штрафується:
		// множимо ще на 1000, адже 1 МВт·год = 1000 кВт·год
		revenueUAH := totalEnergy * portion * cost * 1000.0

		// Штраф (грн) з частини, яка потрапляє під небаланс
		penaltyUAH := totalEnergy * (1.0 - portion) * cost * 1000.0

		// Прибуток у тис. грн
		profitTys := (revenueUAH - penaltyUAH) / 1000.0

		data := InputData{
			PowerMW: powerMW,
			Cost:    cost,
			Portion: portion,

			// Переводимо у тис. грн
			Revenue: revenueUAH / 1000.0,
			Penalty: penaltyUAH / 1000.0,
			Profit:  profitTys,
		}

		// HTML-шаблон із вбудованим CSS
		tmpl, _ := template.New("index").Parse(`
            <!DOCTYPE html>
            <html lang="uk">
            <head>
                <meta charset="UTF-8">
                <title>Контрольний приклад</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        background-color: #f8f9fa;
                        margin: 0;
                        padding: 0;
                    }
                    .container {
                        max-width: 600px;
                        margin: 40px auto;
                        padding: 20px;
                        background-color: #fff;
                        border-radius: 5px;
                        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    }
                    h1, h2 {
                        text-align: center;
                    }
                    form {
                        display: flex;
                        flex-direction: column;
                        margin-bottom: 20px;
                    }
                    label {
                        font-weight: bold;
                        margin: 10px 0 5px;
                    }
                    input[type="text"] {
                        padding: 8px;
                        margin-bottom: 10px;
                        border: 1px solid #ccc;
                        border-radius: 4px;
                    }
                    button[type="submit"] {
                        width: 150px;
                        margin: 0 auto;
                        padding: 10px;
                        background-color: #007bff;
                        color: #fff;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                    }
                    button[type="submit"]:hover {
                        background-color: #0056b3;
                    }
                    .result p {
                        margin: 5px 0;
                    }
                    .note {
                        font-style: italic;
                        color: #777;
                        margin-top: 20px;
                        text-align: center;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <h1>Контрольний приклад</h1>
                    <form action="/" method="post">
                        <label for="powerMW">Потужність (МВт):</label>
                        <input type="text" id="powerMW" name="powerMW" value="{{.PowerMW}}" required>

                        <label for="cost">Ціна (грн/кВт·год):</label>
                        <input type="text" id="cost" name="cost" value="{{.Cost}}" required>

                        <label for="portion">Частка без штрафу (0..1):</label>
                        <input type="text" id="portion" name="portion" value="{{.Portion}}" required>

                        <button type="submit">Розрахувати</button>
                    </form>

                    {{if .Profit}}
                    <div class="result">
                        <h2>Результати</h2>
                        <p>Доходи (тис. грн): {{printf "%.1f" .Revenue}}</p>
                        <p>Штраф (тис. грн): {{printf "%.1f" .Penalty}}</p>
                        <p>Прибуток (тис. грн): {{printf "%.1f" .Profit}}</p>
                    </div>
                    {{end}}
                </div>
            </body>
            </html>
        `)
		tmpl.Execute(w, data)
	} else {
		// Початкова сторінка (порожня форма)
		tmpl, _ := template.New("index").Parse(`
            <!DOCTYPE html>
            <html lang="uk">
            <head>
                <meta charset="UTF-8">
                <title>Контрольний приклад із CSS</title>
                <style>
                    body {
                        font-family: Arial, sans-serif;
                        background-color: #f8f9fa;
                        margin: 0;
                        padding: 0;
                    }
                    .container {
                        max-width: 600px;
                        margin: 40px auto;
                        padding: 20px;
                        background-color: #fff;
                        border-radius: 5px;
                        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    }
                    h1 {
                        text-align: center;
                    }
                    form {
                        display: flex;
                        flex-direction: column;
                        margin-bottom: 20px;
                    }
                    label {
                        font-weight: bold;
                        margin: 10px 0 5px;
                    }
                    input[type="text"] {
                        padding: 8px;
                        margin-bottom: 10px;
                        border: 1px solid #ccc;
                        border-radius: 4px;
                    }
                    button[type="submit"] {
                        width: 150px;
                        margin: 0 auto;
                        padding: 10px;
                        background-color: #007bff;
                        color: #fff;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                    }
                    button[type="submit"]:hover {
                        background-color: #0056b3;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <h1>Контрольний приклад (за методикою з прикладу)</h1>
                    <form action="/" method="post">
                        <label for="powerMW">Потужність (МВт):</label>
                        <input type="text" id="powerMW" name="powerMW" required>

                        <label for="cost">Ціна (грн/кВт·год):</label>
                        <input type="text" id="cost" name="cost" required>

                        <label for="portion">Частка без штрафу (0..1):</label>
                        <input type="text" id="portion" name="portion" required>

                        <button type="submit">Розрахувати</button>
                    </form>
                </div>
            </body>
            </html>
        `)
		tmpl.Execute(w, nil)
	}
}
