(function () {
    function makeid(length) {
        var result = '';
        var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        var charactersLength = characters.length;
        for (var i = 0; i < length; i++) {
            result += characters.charAt(Math.floor(Math.random() * charactersLength));
        }
        return result;
    }

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

        msg = "getting schema for workouts";
        console.log(msg);
        tableau.log(msg);

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
                    cols.push({
                        "id": columns[c].name,
                        "alias": columns[c].name,
                        "dataType": typeConversion.get(columns[c].goType)
                    });


                    var tableSchema = {
                        id: "workouts",
                        alias: "workouts",
                        columns: cols
                    };

                    tableSchemas.push(tableSchema);
                }

                msg = "successfully got schema for workouts.";
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

        tableau.log("DEBUG token: " + tableau.password);

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
                    tableau.log("this is your data for " + t + " with " + tableData.length + " rows.");
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
        tableau.log("random check: " + makeid(6));
        tableau.log("phase: " + tableau.phase);

        tableau.authType = tableau.authTypeEnum.custom;
        tableau.connectionName="Peloton Data Connector";

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

        var accessToken = Cookies.get("peloton_wdc_test");
        tableau.log("access token from cookie is '" + accessToken + "'");
        var hasAuth = (accessToken && accessToken.length > 0) || tableau.password.length > 0;
        updateUIWithAuthState(hasAuth);

        initCallback();

        // If we are not in the data gathering phase, we want to store the token.
        // This allows us to access the token in the data gathering phase.
        if (tableau.phase === tableau.phaseEnum.interactivePhase || tableau.phase === tableau.phaseEnum.authPhase) {
            tableau.log("DEBUG phase " + tableau.phase + " where hasAuth = " + hasAuth + " and accessToken = " + accessToken);
            if (hasAuth) {
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
        var accessToken = Cookies.get("peloton_wdc_test");
        var hasAuth = accessToken && accessToken.length > 0;
        updateUIWithAuthState(hasAuth);

        $("#getcyclingdatalink").click(function() {
            tableau.connectionName = "Peloton Data Connector";
            tableau.submit();
        });
    });

    // An on-click function for the Login to Peloton link.
    // This will redirect the user to a login page.
    function doAuthRedirect() {
        var url = "login"
        window.location.href = url;
    }

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