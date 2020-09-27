$(function () {
    let txtKey = $("#txtKey");
    let txtContent = $("#txtContent");
    let saveMethod = "POST";

    let searchParams = new URLSearchParams(window.location.search);
    let key = searchParams.get('key');
    let originalObject;

    if (!!key) {
        $.ajax({
            url: `/api/mock/key/${key}`,
            success: function (data) {
                saveMethod = "PUT";
                originalObject = JSON.parse(data);

                txtKey.val(originalObject.Key);
                txtKey.prop('disabled', true);

                txtContent.val(originalObject.Content);
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

                originalObject = result;
                resetForm();
            },
            error: function (errMsg) {
                console.log(errMsg);
                $.notify(`Error saving mock.\r\n${errMsg.responseText}`, "error");
            }
        })
    });

    $("#btnReset").on("click", function () {
        resetForm();
    });

    function resetForm() {
        if (saveMethod === "POST") {
            txtKey.val("");
            txtContent.val("");
            txtKey.focus();
        } else if (saveMethod === "PUT") {
            txtContent.val(originalObject.Content);
            txtContent.focus();
        } else {
            // Do nothing
        }
    }
});