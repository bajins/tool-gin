/**
 * @Description:
 * @Author: bajins www.bajins.com
 * @File: index.js
 * @Version: 1.0.0
 * @Time: 2019/9/12 11:29
 * @Project: key-gin
 * @Package:
 * @Software: GoLand
 */


$(function () {
    $.ajax({
        url: "/SystemInfo",
        type: "POST",
        dataType: "json",
        success: function (result) {
            $(".version").text(result.data.Version);
        }
    })
})


/**
 * 重置首页版本或产品
 */
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

/**
 * 获取激活码
 */
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
        download("/getKey", {company: company, app: app, version: version});

    } else {
        $.ajax({
            url: "/getKey",
            type: "POST",
            data: {company: company, app: app, version: version},
            contentType: "application/x-www-form-urlencoded; charset=UTF-8",
            dataType: "json",
            success: function (result) {
                let html = "<div style='width:100%;height:100%;padding:5%;'><p><b>产品：</b>" + app + "</p><hr />";
                if (company == "torchsoft") {
                    html = html + "<p><b>许可证数量：</b>" + version + "</p><hr />";
                } else {
                    html = html + "<p><b>版本：</b>" + version + "</p><hr />";
                }
                html = html + "<p><b>key：</b><pre style='background: black;color:#66FF66;padding:5%;'>" + result.data.key + "</pre></p><hr /></div>";
                if (result.code == 200) {
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
                    layer.msg(result.message, {icon: 5});
                }
            }
        })
    }
}


function xshellDownload() {
    let app = $("#xshell-app").val();
    let version = $("#xshell-version").val();
    //加载层,0代表加载的风格，支持0-2,0.5透明度的白色背景
    let index = layer.load(0, {shade: [0.5,'#fff']});
    $.ajax({
        url: "/getXshellUrl",
        type: "POST",
        data: {app: app, version: version},
        contentType: "application/x-www-form-urlencoded; charset=UTF-8",
        dataType: "json",
        success: function (result) {
            layer.close(index);
            if (result.code == 200) {
                let html = "<div style='width:100%;height:100%;padding:5%;'>" +
                    "<p><b>下载地址：</b></p>" +
                    "<p><a href='" + result.data.url + "' target='_blank'>" + result.data.url + "</a></p>" +
                    "</div>";
                let area_width = "40%";
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
                layer.msg(result.message, {icon: 5});
            }
        }
    })
}