$(function () {
    let txtKey = $("#txtKey");
    let txtContent = $("#txtContent");
    let saveMethod = "POST";

    let searchParams = new URLSearchParams(window.location.search);
    let key = searchParams.get('key');

    if (!!key) {
        $.ajax({
            url: `/api/mock/key/${key}`,
            success: function (data) {
                saveMethod = "PUT";
                let jObject = JSON.parse(data);

                txtKey.val(jObject.Key);
                txtKey.prop('disabled', true);

                txtContent.val(jObject.Content);
            },
            error: function (errMsg) {
                $.notify(`Mock with Key '${key}' is not existing but you can create it now.`, "warning");
                txtKey.val(key);
                txtContent.focus();
            }
        });
    }

    $("#btnSubmit").on("click", function () {
        $.ajax({
            type: saveMethod,
            url: '/api/mock',
            data: JSON.stringify({
                "Key": txtKey.val(),
                "Content": txtContent.val()
            }),
            success: function (data) {
                const result = JSON.parse(data);
                $.notify(`Mock '${result.Key}' saved.`, "success");

                if (saveMethod === "POST") {
                    txtKey.val("");
                    txtContent.val("");
                    txtKey.focus();
                }
            },
            error: function (errMsg) {
                console.log(errMsg);
                $.notify(`Error saving mock.\r\n${errMsg.responseText}`, "error");
            }
        })
    });
});