(function () {
    var myConnector = tableau.makeConnector();

    myConnector.getSchema = function (schemaCallback) {
        let typeConversion = new Map();
        typeConversion.set("string", "tableau.dataTypeEnum.string");
        typeConversion.set("int32", "tableau.dataTypeEnum.int");
        typeConversion.set("uint16", "tableau.dataType.int");
        typeConversion.set("uint64", "tableau.dataType.int");
        typeConversion.set("float32", "tableau.dataTypeEnum.float");
        typeConversion.set("bool", "tableau.dataTypeEnum.bool");
        typeConversion.set("time.Time", "tableau.dataTypeEnum.datetime");

        var tableSchemas = [];

        tableau.log("This is your Peloton schema.");
        $.ajax({
            // beforeSend: function (request) {
            //     request.setRequestHeader("Access-Control-Allow-Origin", "*");
            // },
            dataType: "json",
            url: "http://localhost:30000/cycling/schema",
        }).done(function (data) {
            var tables = data.tables
            for (var t = 0, tlen = tables.length; t < tlen; t++) {
                /*
                "name", "description", "columns"
                */
                var cols = []
                table = tables[t]
                columns = table.columns
                for (var c = 0, clen = columns.length; c < clen; c++) {
                    cols.push({
                        "id": columns[c].name,
                        "alias": columns[c].name,
                        "dataType": typeConversion.get(columns[c].goType)
                    });
                }

                var tableSchema = {
                    id: table.name,
                    alias: table.name,
                    columns: cols
                };

                tableSchemas.push(tableSchema);
            }

            schemaCallback(tableSchemas);
        });
    };

    myConnector.getData = function (table, doneCallback) {
        doneCallback();
    };

    tableau.registerConnector(myConnector);

    $(document).ready(function () {
        $("#submitButton").click(function () {
            tableau.connectionName = "Peloton Data Feed";
            tableau.submit();
        });
    });
})();