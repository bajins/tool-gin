function selectCompany() {
    let company = $("#company").val();
    if (company == "netsarang") {
        $("#app").empty();
        //为Select追加一个Option(下拉项)
        $("#app").append('<option value="Xmanager">Xmanager Power Suite</option>');
        $("#app").append('<option value="Xshell">Xshell</option>');
        $("#app").append('<option value="Xlpd">Xlpd</option>');
        $("#app").append('<option value="Xftp">Xftp</option>');
        $("#app").append('<option value="Xshell Plus" selected>Xshell Plus</option>');

        $("#version-label").text("版本:");
        $("#version").empty();
        $("#version").append('<option value="6" selected>6</option>');
        $("#version").append('<option value="5">5</option>');
        $("#version").append('<option value="4">4</option>');
        $("#version").append('<option value="3">3</option>');
        $("#version").append('<option value="2">2</option>');
    } else if (company == "mobatek") {
        $("#app").empty();
        //为Select追加一个Option(下拉项)
        $("#app").append('<option value="MobaXterm" selected>MobaXterm</option>');

        $("#version-label").text("版本:");
        $("#version").empty();
        $("#version").append('<option value="11.1" selected>11.1</option>');
        $("#version").append('<option value="11.0">11.0</option>');
        $("#version").append('<option value="10.9">10.9</option>');
        $("#version").append('<option value="10.8">10.8</option>');
        $("#version").append('<option value="10.7">10.7</option>');
        $("#version").append('<option value="10.6">10.6</option>');
        $("#version").append('<option value="10.5">10.5</option>');
        $("#version").append('<option value="10.4">10.4</option>');
        $("#version").append('<option value="10.2">10.2</option>');
        $("#version").append('<option value="10.1">10.1</option>');
        $("#version").append('<option value="10.0">10.0</option>');
        $("#version").append('<option value="9.4">9.4</option>');
        $("#version").append('<option value="9.3">9.3</option>');
        $("#version").append('<option value="9.2">9.2</option>');
        $("#version").append('<option value="9.1">9.1</option>');
        $("#version").append('<option value="9.0">9.0</option>');
    } else if (company == "torchsoft") {
        $("#app").empty();
        //为Select追加一个Option(下拉项)
        $("#app").append('<option value="Registry Workshop" selected>Registry Workshop</option>');

        $("#version-label").text("许可证数量:");
        $("#version").empty();
        $("#version").append('<option value="10">10</option>');
        $("#version").append('<option value="9">9</option>');
        $("#version").append('<option value="8">8</option>');
        $("#version").append('<option value="7">7</option>');
        $("#version").append('<option value="6">6</option>');
        $("#version").append('<option value="4">4</option>');
        $("#version").append('<option value="3">3</option>');
        $("#version").append('<option value="2">2</option>');
        $("#version").append('<option value="1" selected>1</option>');
    }
}


function getKey() {
    let company = $("#company").val();
    let app = $("#app").val();
    let version = $("#version").val();
    if (app == "MobaXterm") {
        // 构造隐藏的form表单
        /*let form = $('<form action="/getKey" method="post">' +
            '<input type="text" name="company" value="' + company + '"/>' +
            '<input type="text" name="app" value="' + app + '"/>' +
            '<input type="text" name="version" value="' + version + '"/>' +
            '</form>');
        $(document.body).append(form);
        form.submit().remove();*/

        $.ajax({
            url: "/getKey",
            type: "POST",
            data: {company: company, app: app, version: version},
            contentType: "application/x-www-form-urlencoded; charset=UTF-8",
            responseType: "blob",
            success: function (result, status, xhr) {
                // 从response的headers中获取filename, 后端response.setHeader("Content-Disposition", "attachment; filename=xxxx.xxx") 设置的文件名;
                let contentDisposition = xhr.getResponseHeader('Content-Disposition');
                let patt = new RegExp("filename=([^;]+\\.[^\\.;]+);*");
                let filename = "";
                // 如果从Content-Disposition中取到的文件名为空
                if (isEmpty(contentDisposition)) {
                    let f = xhr.config.params.filePath.split("/");
                    filename = f[f.length - 1];
                } else {
                    filename = patt.exec(contentDisposition)[1];
                }
                // 取文件名信息中的文件名,替换掉文件名中多余的符号
                filename = filename.replaceAll("\\\\|/|\"", "");

                let downloadElement = document.createElement('a');
                downloadElement.style.display = 'none';

                //这里res.data是返回的blob对象
                let blob = new Blob([result], {type: 'application/octet-stream;charset=utf-8'});
                // 创建下载的链接
                let href = window.URL.createObjectURL(blob);
                downloadElement.href = href;
                // 下载后文件名
                downloadElement.download = filename;
                document.body.appendChild(downloadElement);
                // 点击下载
                downloadElement.click();
                // 下载完成移除元素
                document.body.removeChild(downloadElement);
                // 释放掉blob对象
                window.URL.revokeObjectURL(href);
            }
        })
    } else {
        $.ajax({
            url: "/getKey",
            type: "POST",
            data: {company: company, app: app, version: version},
            contentType: "application/x-www-form-urlencoded; charset=UTF-8",
            dataType: "text",
            success: function (result) {
                let res = JSON.parse(result);
                let html = "<div style='width:100%;height:100%;padding:5%;'><p><b>产品：</b>" + app + "</p><hr />";
                if (company == "torchsoft") {
                    html = html + "<p><b>许可证数量：</b>" + version + "</p><hr />";
                } else {
                    html = html + "<p><b>版本：</b>" + version + "</p><hr />";
                }
                html = html + "<p><b>key：</b><pre style='background: black;color:#66FF66;padding:5%;'>" + res.data.key + "</pre></p><hr /></div>";
                if (res.code == 200) {
                    let area_width = "30%";
                    if (device.isMobile) {
                        area_width = "80%";
                    }
                    //自定页
                    layer.open({
                        // 在默认状态下，layer是宽高都自适应的，但当你只想定义宽度时，你可以area: '500px'，高度仍然是自适应的。
                        // 当你宽高都要定义时，你可以area: ['500px', '300px']
                        area: [area_width],
                        type: 1,
                        icon: 1,
                        skin: 'layui-layer-lan', //样式类名,目前layer内置的skin有：layui-layer-lan、layui-layer-molv
                        closeBtn: 1, //关闭按钮
                        anim: 2,
                        shadeClose: true, //开启遮罩关闭
                        title: false,
                        content: html
                    });
                } else {
                    //提示层
                    layer.msg(res.message, {icon: 5});
                }
            }
        })
    }
}

