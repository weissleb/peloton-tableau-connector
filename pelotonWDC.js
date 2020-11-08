(function () {
    var myConnector = tableau.makeConnector();

    myConnector.getSchema = function (schemaCallback) {
        let typeConversion = new Map();
        typeConversion.set("string", tableau.dataTypeEnum.string);
        typeConversion.set("int32", tableau.dataTypeEnum.int);
        typeConversion.set("uint16", tableau.dataTypeEnum.int);
        typeConversion.set("uint64", tableau.dataTypeEnum.int);
        typeConversion.set("float32", tableau.dataTypeEnum.float);
        typeConversion.set("bool", tableau.dataTypeEnum.bool);
        typeConversion.set("time.Time", tableau.dataTypeEnum.datetime);

        var tableSchemas = [];

        $.getJSON("http://localhost:30000/cycling/schema?tables=workouts", function(resp) {
            var tables = resp.tables
            for (var t = 0, tlen = tables.length; t < tlen; t++) {
                // "name", "description", "columns"
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

            tableau.log("This is your Peloton schema.");
            schemaCallback(tableSchemas);
        });
    };

    myConnector.getData = function (table, doneCallback) {
        $.getJSON("http://localhost:30000/cycling/data/" + table.tableInfo.id, function(resp) {
            var tableData = [];

            if (table.tableInfo.id == "workouts") {
                tableData = resp.data;
                tableau.log("This is your Peloton data with " + tableData.length + " rows.");
            }

            if (table.tableInfo.id == "dummy") {
                row = {
                    "stringField": "foo",
                    "intField": 99,
                    "boolField": true,
                    "floatField": 99.9,
                    "dateField": "1976-01-07 11:05:30"
                };
                tableData.push(row)
            }

            table.appendRows(tableData)
            doneCallback();
        });
    };

    tableau.registerConnector(myConnector);

    $(document).ready(function () {
        $("#submitButton").click(function () {
            tableau.connectionName = "Peloton Data Feed";
            tableau.submit();
        });
    });
})();

function f() {
    console.log("Fuck you!")
}