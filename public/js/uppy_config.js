const { ProgressBar } = Uppy;
const { XHRUpload } = Uppy;
const { AwsS3 } = Uppy;
const { Dashboard } = Uppy;
const { ImageEditor } = Uppy;

baseUrl = "http://localhost:8080/";

const getS3PreSignUrl = (file) => {
    var url = baseUrl + "get-presigned?name=" + file.fileName + "&type=" + file.type;
    data = $.get(url);
    return data;
}

const createTable = () => {
	var url = baseUrl + "files";
	$('#table-body').empty();
	table.DataTable({
		"pageLength": 50,
		"lengthMenu": [10, 25, 50, 100, 1000],
		"processing": true,
		"ordering": false,
		"paging": true,
		"serverSide": true,
		"ajax": {
			url: url, type: "post"
		},
		"columns": [
			{ "data": "name", "targets": 0, },
			{ "data": "owner", "targets": 1, },
			{
				"data": "id", "render": function (data, type, row) {
					return '<button class="file-delete" data-id=' + row.id + '>Delete</button>';
				}, "targets": 2,
			},
		]
	});
	$('.dataTables_length').addClass('bs-select');
	new ClipboardJS('.copyToClipboard');
};

const refreshTable = () => {
	var tempTable = $('#dataTable').DataTable();
	tempTable.destroy();
	createTable();
};

const onUploadSuccess = (elForUploadedFiles) => async (file, response) => {
    let fileName = file.name;
    let fileType = file.type;
    let fileSize = file.size;
    let url = baseUrl + "files?name=" + fileName;
    let li = document.createElement('li');
    let a = document.createElement('a');
    let copyDiv = document.createElement('div');
    let icon = document.createElement('i');

    icon.className = 'fas fa-copy copyToClipboard';
    icon.setAttribute("data-clipboard-text", url);
    copyDiv.className = "tooltipped tooltipped - n m - 2 p - 2 recentUpload";
    copyDiv.setAttribute('aria-label', 'Copy to Clipboard');
    copyDiv.appendChild(icon);

    a.href = url;
    a.target = '_blank';
    a.appendChild(document.createTextNode(fileName));
    li.appendChild(a);
    li.appendChild(copyDiv);

    document.querySelector(elForUploadedFiles).appendChild(li);
    let urlChunks = response.uploadURL.split("/");
    let fullFileName = urlChunks[urlChunks.length - 1];
    await $.get(baseUrl + "store-file?name=" + fullFileName + "&contentType=" + fileType + "&size=" + fileSize);

    refreshTable();
};

const restrictionError = (elForErrorMessage) => (file, error) => {
    let fileName = file.name;
    var errorMessage = error.toString().replace(fileName, "");
    errorMessage = errorMessage.replace("Error:", "Error (" + fileName + "):");
    let li = document.createElement('li');
    li.appendChild(document.createTextNode(errorMessage));
    document.querySelector(elForErrorMessage).appendChild(li);
};

const uploadError = (elForErrorMessage) => (file, error, response) => {
    let fileName = file.name;
    let errorMessage = fileName + ": " + error;
    let li = document.createElement('li');
    li.appendChild(document.createTextNode(errorMessage));
    document.querySelector(elForErrorMessage).appendChild(li);
};

const consoleError = (elForErrorMessage) => (error) => {
    console.log(error);
    let errorMessage = "There was an error, please check the console.";
    let li = document.createElement('li');
    li.appendChild(document.createTextNode(errorMessage));
    document.querySelector(elForErrorMessage).appendChild(li);
};

$(document).ready(function () {
    const uppy = new Uppy.Core({
        debug: true, autoProceed: false,
        restrictions: {
            maxFileSize: 2147483648, // 2gb
            maxNumberOfFiles: 15,
            //allowedFileTypes: ['.jpg', '.jpeg', '.png', '.gif', '.pdf', '.csv', '.xslx', '.txt', '.js', '.css'],
        },
    });
    uppy
        .use(Dashboard, {
            target: '.upload-form .for-DragDrop',
            inline: true,
            showProgressDetails: true,
        })
        .use(AwsS3, {
            getUploadParameters: file => {
                return getS3PreSignUrl({
                    fileName: file.name,
                    type: file.type
                })
                    .then(data => {
                        return {
                            method: 'PUT',
                            url: data.url,
                            fields: {},
                            headers: {
                                'content-type': data.type
                            }
                        }
                    })
            }
        })
        .use(ImageEditor, { target: Dashboard })
        .on('upload-success', onUploadSuccess('.upload-form .uploaded-files ol'))
        .on('upload-error', restrictionError('.upload-form .error-messages ol'))
        .on('error', consoleError('.upload-form .error-messages ol'))
        .on('restriction-failed', restrictionError('.upload-form .error-messages ol'));

        createTable();

});