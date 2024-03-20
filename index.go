package process_dashboard

func GetIndexContent() string {
	return `
<html>
    <head>
        <title>Dashboard</title>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.4/Chart.js"></script>
        <link rel="preconnect" href="https://fonts.googleapis.com">
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
        <link href="https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap" rel="stylesheet">

        <style>
            body {
                background-color: #F2F3F4;
                color: #17A589;
                font-family: "Roboto", sans-serif;
            }
            .title {
                font-size: large;
                margin-bottom: 1em;
                font-weight: bold;
            }
            .card {
                height: 30%;
                max-height: 350px;
                margin: 1em;
                margin-right: 1em;
                margin-left: 1em;
                background-color: #fff;
                padding: 1em;
            }
        </style>
    </head>
    <body>
        [CANVAS]
        <script>
            var data = [];

            function doStuff() {
                [CODE]
            }

            async function fetchData() {
                const res = await fetch("/data");
                data = await res.json();
                doStuff();
            }

            setInterval(fetchData, 3000)
        </script>
    </body>
</html>`
}
