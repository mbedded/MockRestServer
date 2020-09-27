// Global variables
var tableContent;

$(function () {
    tableContent = $("#tableContent");

    $.ajax({
        type: "GET",
        url: '/api/mock/all',
        success: function (data) {
            const result = JSON.parse(data);

            createTable(result);
        },
        error: function (errMsg) {
            $.notify(`Error reading all mocks.\r\n${errMsg.responseText}`, "error");
        }
    });
});

function createTable(mocks) {
    let rows = [mocks.length];

    if (mocks.length <= 0) {
        rows[0] = `<tr>
                    <td class='text-muted' colspan='3'>
                        There are no Mocks existing in the database but you can <a href="/create">create a new one</a>.
                    </td>
                   </tr>`
    } else {
        for (let index = 0; index < mocks.length; index++) {
            let item = mocks[index];

            let key = item.Key;

            let addition = (item.Content.length <= 20) ? "" : "…";
            let content = item.Content.substring(0, 20) + addition;
            let rowContent = `
                            <tr id=\"row_${key}\">
                                <td>${key}</td>
                                <td>${content}</td>
                                <td>
                                    <a href=\"/raw/${key}\" target="_blank">Show raw content</a> |
                                    <a href=\"/create?key=${key}\">Edit</a> |
                                    <a class="text-danger" href="#" onclick=\"deleteMock(\'${key}\')\">Delete</a>
                                </td>
                            </tr>`;

            rows[index] = rowContent;
        }
    }

    var complete = rows.join();
    tableContent.html(complete);
}

function deleteMock(key) {
    let shouldBeDeleted = confirm(`Do you want to delete: ${key}`)
    if (shouldBeDeleted === false) {
        return;
    }

    $.ajax({
        type: "DELETE",
        url: `/api/mock/key/${key}`,
        success: function () {
            $.notify(`Mock '${key}' has been deleted.`, "success");

            let rowId = `#row_${key}`;
            $(rowId).remove();
        },
        error: function (errMsg) {
            $.notify(`Error deleting Mock.\r\n${errMsg.responseText}`, "error");
        }
    });
}
