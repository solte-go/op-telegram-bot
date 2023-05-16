function requestProjectsStatus() {
    fetch("/api/v1/healthCheck", {
        method: 'GET',
        headers: {'Content-Type': 'application/json'},
    })
        .then((response) => {
            return response.json();
        })
        .then((data) => {
            const myNode = document.getElementById("projects-container");
            while (myNode.firstChild) {
                myNode.removeChild(myNode.lastChild);
            }

            data.map(function (project) {
                let syncDate = new Date(`${project.LastSync}`).toLocaleString();
                let errorField;
                let projectID = `${project.Name}`;
                let projectSyncStatus;
                switch(project.Status) {
                    case "stopped":
                        projectSyncStatus = "row-process-status-paused"
                        break;
                    case "ok":
                        projectSyncStatus = "row-process-status-ok"
                        break;
                    case "error":
                        projectSyncStatus = "row-process-status-error"
                        break;
                    default:
                        projectSyncStatus = "row-process-status-paused"
                }

                let error = `${project.Status}`;
                if (error !== "ok") {
                    errorField = "header-title-status-error"
                } else {
                    errorField = "header-title-status-ok"
                }
                let projectName, projectStatus, projectSyncDate;

                projectName =  "<span class='header-title-name'> Project: </span>" + project.Name + " / ";

                projectStatus = "<span class='header-title-name'> Status: </span>" +
                    "<span class="+errorField+">" + project.Status + "</span> / ";

                projectSyncDate = "<span class='header-title-name'>Last sync: </span>" +syncDate;

                let statusLine = (projectName + projectStatus + projectSyncDate);

                $('#projects-container').append([
                    $('<div>', {"class": "project-row"}).append([
                        $('<div>', {"class": "project-header"}).append([
                            $('<span>', {"class":"m-2 project-title-bar",
                                // "type": "button",
                                "data-bs-toggle":"collapse",
                                "data-bs-target":"#"+ projectID.replace(/\s/g,''),
                                "aria-expanded":"false",
                                "aria-controls":project.Name}).val(project.Name).html(statusLine),
                            $('<span>',{"class":"float-start " + projectSyncStatus}),
                            $('<button>', {"class":"btn btn-sm btn-danger mt-2 me-1 ms-1 float-end",
                                "type": "button",
                                "onclick":'getProjectDataForDeletion(this.value)'}).val(project.Name).text("Delete"),
                            $('<button>', {"class":"btn btn-sm btn-success mt-2 float-end",
                                "type": "button",
                                "onclick":'getProjectDataForUpdate(this.value)'}).val(project.Name).text("Edit"),
                            $('<button>', {"class":"btn btn-sm btn-secondary mt-2 me-1 float-end",
                                "type": "button",
                                "onclick":'changeProjectSyncOption(this.value)'}).val(project.Name).text("Sync"),
                        ]),
                        $('<div>', {"class":"collapse project-details", "id":projectID.replace(/\s/g,'')}).append([
                            $('<div>', {"class":"row projects-grid" }).append([
                                $('<div>', {"class":"col-md-auto projects-grid-col-names" }).append([
                                    $('<span>',{"class":"badge bg-secondary badge-cm-style p-2 me-1"}).text("Project status"),
                                ]),
                                $('<div>', {"class":"col" }).append([
                                    $('<span>',{"class":"project-content-font"}).text(project.Status),
                                ]),
                            ]),
                            $('<div>', {"class":"row projects-grid" }).append([
                                $('<div>', {"class":"col-md-auto projects-grid-col-names" }).append([
                                    $('<span>',{"class":"badge bg-secondary badge-cm-style p-2 me-1"}).text("Project last sync date"),
                                ]),
                                $('<div>', {"class":"col" }).append([
                                    $('<span>',{"class":"project-content-font"}).text(syncDate),
                                ]),
                            ]),
                            $('<div>', {"class":"row projects-grid" }).append([
                                $('<div>', {"class":"col-md-auto projects-grid-col-names" }).append([
                                    $('<span>',{"class":"badge bg-secondary badge-cm-style p-2 me-1"}).text("Project last error"),
                                ]),
                                $('<div>', {"class":"col" }).append([
                                    $('<span>',{"class":"project-content-font"}).text(project.LastSyncError),
                                ]),
                            ]),
                            $('<div>', {"class":"row projects-grid" }).append([
                                $('<div>', {"class":"col-md-auto projects-grid-col-names" }).append([
                                    $('<span>',{"class":"badge bg-secondary badge-cm-style p-2 me-1"}).text("Project manager"),
                                ]),
                                $('<div>', {"class":"col" }).append([
                                    $('<span>',{"class":"project-content-font"}).text(project.Manager),
                                ]),
                            ]),
                            $('<div>', {"class":"row projects-grid" }).append([
                                $('<div>', {"class":"col-md-auto projects-grid-col-names" }).append([
                                    $('<span>',{"class":"badge bg-secondary badge-cm-style p-2 me-1"}).text("Project query"),
                                ]),
                                $('<div>', {"class":"col" }).append([
                                    $('<span>',{"class":"project-content-font"}).text(project.Query),
                                ]),
                            ]),
                            $('<div>', {"class":"row" }).append([
                                $('<div>', {"class":"col project-sla-content" }).append([
                                    $('<div>').append([
                                        $('<h4>', {"class":"project-content-font-weight"}).text("SLA").css("margin-bottom","0px"),
                                    ]),
                                    $('<span>',{"class":"badge bg-secondary p-2 me-1"}).text("Time To React"),
                                    $('<span>',{"class":"project-content-font-weight"}).text(project.SLA.TimeToReact),
                                    $('<span>',{"class":"badge bg-secondary p-2 me-1"}).text("Time To Resolve"),
                                    $('<span>',{"class":"project-content-font-weight"}).text(project.SLA.TimeToResolve),
                                    $('<span>',{"class":"badge bg-secondary p-2 me-1"}).text("Stuck Issue"),
                                    $('<span>',{"class":"project-content-font-weight"}).text(project.SLA.IssueStuck),
                                ]),
                            ]),
                        ]),
                    ]),
                ]);
            });
        })
        .catch(function (error) {
            console.log(error);
        });
}

function requestProjectsNames() {
    fetch("/api/v1/projects", {
        method: 'GET',
        headers: {'Content-Type': 'application/json'},
    })
        .then((response) => {
            return response.json();
        })
        .then((data) => {
            const myNode = document.getElementById("projects");
            while (myNode.firstChild) {
                myNode.removeChild(myNode.lastChild);
            }

            data.map(function (service) {
                $('#projects').append([
                    $('<option>',{"value": service}).text(service)
                ]);
            });
        })
        .catch(function (error) {
            console.log(error);
        });
}

function userLogOut(){
     fetch("/api/authentication/logout", {
         method: "POST",
         credentials: "include",
         headers: {'Content-Type': 'application/json'},
    })
         .then((response) => {
             if (!response.ok) {
                 $("#LoginModal").modal('hide');
                 $("#page_alerts").css("visibility", "unset");
                 $("#page_alerts_text").text("Please login");
             }
         });
    location.reload();
}

function createUser(){
    let email = "";
    let name = "";
    let password = "";
    let obj = {
        email: email,
        name: name,
        password:password,
    };

    fetch("/api/v1/createUser", {
        method: "POST",
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(obj),
    })
        .then((response) => {
            if (!response.ok) {
                $("#page_alerts").css("visibility", "unset");
                $("#page_alerts_text").text(myJson.error);
            } else {
                console.log(JSON.stringify(myJson));
                location.reload();
        }
        return response.json();
        })
        .catch(function (error) {
            console.log(error);
        });
}

function editProjectsRequest() {
    fetch("/api/v1/projects", {
        method: 'GET',
        headers: {'Content-Type': 'application/json'},
    })
        .then((response) => {
            return response.json();
        })
        .then((data) => {
            const myNode = document.getElementById("editProjectName");
            while (myNode.firstChild) {
                myNode.removeChild(myNode.lastChild);
            }

            data.map(function (service) {
                $('#editProjectName').append([
                    $('<option>',{"value": service}).text(service)
                ]);
            });
        })
        .catch(function (error) {
            console.log(error);
        });
}

function reloadPage(){
    location.reload();
}