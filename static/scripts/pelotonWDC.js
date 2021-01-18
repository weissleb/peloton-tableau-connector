(function () {

    var myConnector = tableau.makeConnector();

    myConnector.getSchema = function (schemaCallback) {
        let typeConversion = new Map();
        typeConversion.set("string", tableau.dataTypeEnum.string);
        typeConversion.set("int", tableau.dataTypeEnum.int);
        typeConversion.set("int32", tableau.dataTypeEnum.int);
        typeConversion.set("uint16", tableau.dataTypeEnum.int);
        typeConversion.set("uint64", tableau.dataTypeEnum.int);
        typeConversion.set("float32", tableau.dataTypeEnum.float);
        typeConversion.set("float64", tableau.dataTypeEnum.float);
        typeConversion.set("bool", tableau.dataTypeEnum.bool);
        typeConversion.set("time.Time", tableau.dataTypeEnum.datetime);

        var tableSchemas = [];

        tableau.log("getting schema for workouts");

        var xhr = $.ajax({
            url: "cycling/schema/workouts",
            type: "GET",
            dataType: 'json',
            async: false,
            success: function (data) {
                var table = data.tables[0];
                // "name", "description", "columns"
                var cols = [];
                columns = table.columns;

                for (var c = 0, clen = columns.length; c < clen; c++) {

                    var columnRole;
                    if (columns[c].name == "RideLengthMinutes" || columns[c].name == "RideDifficulty" || columns[c].name == "StartTimeSeconds") {
                        columnRole = tableau.columnRoleEnum.dimension.valueOf();
                    } else {
                        columnRole = undefined;
                    }

                    var columnType;
                    if (columns[c].name == "RideLengthMinutes" || columns[c].name == "RideDifficulty") {
                        columnType = tableau.columnTypeEnum.discrete.valueOf();
                    } else {
                        columnType = undefined;
                    }

                    var numberFormat;
                    if (columns[c].name == "AvgResistance") {
                        numberFormat = tableau.numberFormatEnum.percentage.valueOf();
                    } else {
                        numberFormat = undefined;
                    }

                    var aggType;
                    if (columns[c].name.startsWith("Avg")) {
                        aggType = tableau.aggTypeEnum.avg.valueOf();
                    } else {
                        aggType = undefined;
                    }

                    tableau.log("DEBUG: column " + columns[c].name + " has role of " + columnRole + ", type of " + columnType +
                        " numberFormat of " + numberFormat + " and aggType of " + aggType);
                    cols.push({
                        "id": columns[c].name,
                        "alias": columns[c].name,
                        "dataType": typeConversion.get(columns[c].goType),
                        "columnRole": columnRole,
                        "columnType": columnType,
                        "numberFormat": numberFormat,
                        "aggType": aggType
                    });


                    var tableSchema = {
                        id: "workouts",
                        alias: "Workouts",
                        description: "Cycling workout with summary metrics.",
                        columns: cols
                    };

                    tableSchemas.push(tableSchema);
                }

                msg = "successfully got schema for workouts";
                console.log(msg);
                tableau.log(msg);
                schemaCallback(tableSchemas);
            },
            error: function (xhr, ajaxOptions, thrownError) {
                tableau.log(xhr.responseText + "\n" + thrownError);
                tableau.abortWithError("error getting schema for workouts");
            }
        });
    };

    myConnector.getData = function (table, doneCallback) {
        tableau.log("trying to get data");
        if (tableau.password.length === 0) {
            tableau.log("we do not have a token, aborting for auth");
            tableau.abortForAuth();
        }

        tableau.log("got token with length of " + tableau.password.length);

        t = table.tableInfo.id;
        tableau.log("getting data for " + t);

        var xhr = $.ajax({
            url: "cycling/data/" + t,
            type: "GET",
            dataType: 'json',
            async: false,
            headers: {
                "Authorization": "Bearer " + tableau.password
            },
            success: function (data) {
                var tableData = [];

                if (t === "workouts") {
                    tableData = data.data;
                    tableau.log("this is your data for " + t + " with " + tableData.length + " rows");
                }

                table.appendRows(tableData);
                doneCallback();
            },
            error: function (xhr, ajaxOptions, thrownError) {
                tableau.log(xhr.responseText + "\n" + thrownError);
                tableau.abortWithError("error getting data for " + t);
            }
        });
    };

    // Init function for connector, called during every phase but
    // only called when running inside the simulator or tableau.
    myConnector.init = function (initCallback) {
        tableau.log("phase: " + tableau.phase);

        tableau.authType = tableau.authTypeEnum.custom;
        // tableau.connectionName="Peloton Data Connector";

        if (tableau.phase === tableau.phaseEnum.gatherDataPhase) {

            // If we don't have a valid token stored in password, we need to
            // re-authenticate.
            if (tableau.password === 0) {
                tableau.log("gatherDataPhase abortForAuth()");
                tableau.abortForAuth();
            } else {
                tableau.log("gatherDataPhase has password and proceed ahead");
            }
        }

        var accessToken = Cookies.get("peloton_wdc_token");
        var user = Cookies.get("peloton_wdc_user");
        if (accessToken && accessToken.length > 0) {
            tableau.log("found cookie with access token of length " + accessToken.length);
        } else {
            tableau.log("did not find cookie or token length is not > 0");
        }
        var hasAuth = (accessToken && accessToken.length > 0) || tableau.password.length > 0;
        updateUIWithAuthState(hasAuth);

        initCallback();

        // If we are not in the data gathering phase, we want to store the token.
        // This allows us to access the token in the data gathering phase.
        if (tableau.phase === tableau.phaseEnum.interactivePhase || tableau.phase === tableau.phaseEnum.authPhase) {
            tableau.log("phase " + tableau.phase + " where hasAuth = " + hasAuth + " and accessToken with length = " + accessToken.length);
            if (hasAuth) {
                tableau.username = user;
                tableau.password = accessToken;

                if (tableau.phase === tableau.phaseEnum.authPhase) {
                    // Auto-submit here if we are in the auth phase
                    tableau.submit()
                }

                return;
            }
        }
    };

    tableau.registerConnector(myConnector);

    $(document).ready(function () {
        var accessToken = Cookies.get("peloton_wdc_token");
        var user = Cookies.get("peloton_wdc_user");
        var hasAuth = accessToken && accessToken.length > 0;
        updateUIWithAuthState(hasAuth);

        $("#getcyclingdatalink").click(function () {
            tableau.connectionName = "Peloton Data Connector for " + user;
            tableau.submit();
        });
    });

    // This function toggles the label shown depending
    // on whether or not the user has been authenticated.
    function updateUIWithAuthState(hasAuth) {
        if (hasAuth) {
            $(".notsignedin").css('display', 'none');
            $(".signedin").css('display', 'block');
        } else {
            $(".notsignedin").css('display', 'block');
            $(".signedin").css('display', 'none');
        }
    }

})();