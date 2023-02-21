/**
 * @Description:
 * @Author: bajins www.bajins.com
 * @File: index.js
 * @Version: 1.0.0
 * @Time: 2019/9/12 11:29
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
import log from './utils/log.js';
import util from './utils/util.js';
import http from "./utils/http.js";
import array from "./utils/array.js";
import string from "./utils/string.js";


/**
 * 监听窗口变化
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
    window.onload = function () {

    }
    window.onresize = function () {

    }
})


// ==================================  获取Netsarang激活key  ===================================

/**
 * 重置首页版本或产品
 */
window.selectCompany = function selectCompany() {
    let company = $("#company").val();
    if (company == "netsarang") {
        //为Select追加一个Option(下拉项)
        $("#app").empty().append(`<option value="Xmanager">Xmanager Power Suite</option>
<option value="Xshell">Xshell</option>
<option value="Xlpd">Xlpd</option>
<option value="Xftp">Xftp</option>
<option value="Xshell Plus" selected>Xshell Plus</option>`);
        $("#version-label").text("版本:");
        $("#version").empty().append(`<option value="7" selected>7</option>
<option value="6">6</option>
<option value="5">5</option>
<option value="4">4</option>
<option value="3">3</option>
<option value="2">2</option>`);
        $("#netSarangDownloadBtn").show();
    } else if (company == "mobatek") {
        $("#app").empty();
        //为Select追加一个Option(下拉项)
        $("#app").append('<option value="MobaXterm" selected>MobaXterm</option>');

        $("#version-label").text("版本:");
        $("#version").empty().append(`<option value="21.0" selected>21.0</option>
<option value="20.0">20.0</option>
<option value="12.4">12.4</option>
<option value="12">12.0</option>
<option value="11.1">11.1</option>
<option value="11.0">11.0</option>
<option value="10.9">10.9</option>
<option value="10.8">10.8</option>
<option value="10.7">10.7</option>
<option value="10.6">10.6</option>
<option value="10.5">10.5</option>
<option value="10.4">10.4</option>
<option value="10.2">10.2</option>
<option value="10.1">10.1</option>
<option value="10.0">10.0</option>
<option value="9.4">9.4</option>
<option value="9.3">9.3</option>
<option value="9.2">9.2</option>
<option value="9.1">9.1</option>
<option value="9.0">9.0</option>`);

        $("#netSarangDownloadBtn").hide();
    } else if (company == "torchsoft") {
        //为Select追加一个Option(下拉项)
        $("#app").empty().append('<option value="Registry Workshop" selected>Registry Workshop</option>');

        $("#version-label").text("许可证数量:");
        $("#version").empty().append(`<option value="10">10</option>
<option value="9">9</option>
<option value="8">8</option>
<option value="7">7</option>
<option value="6">6</option>
<option value="4">4</option>
<option value="3">3</option>
<option value="2">2</option>
<option value="1" selected>1</option>`);
        $("#netSarangDownloadBtn").hide();
    }
}

/**
 * 获取激活码
 *
 * ECMAScript6使用全局变量配置页面绑定事件
 *
 */
window.getKey = function getKey() {
    let company = $("#company").val();
    let app = $("#app").val();
    let version = $("#version").val();
    if (app == "MobaXterm") {
        // 构造隐藏的form表单
        /*let form = $(`<form action="/getKey" method="post" target="_blank">
            <input type="text" name="company" value="${company}"/>
            <input type="text" name="app" value="${app}"/>
            <input type="text" name="version" value="${version}"/>
            </form>`);
        $(document.body).append(form);
        form.submit().remove();*/
        http.ajax({
            url: "/getKey",
            method: http.METHOD.POST,
            data: {company: company, app: app, version: version},
            responseType: http.RESPONSE_TYPE.BLOB,
            success: (result, status, xhr) => {
                log.info(result);
            },
            error: (xhr, status, error) => {
                console.log(error)
                layer.msg(error, {icon: 5});
            }
        });
    } else {
        $.ajax({
            url: "/getKey",
            type: "POST",
            data: {company: company, app: app, version: version},
            contentType: "application/x-www-form-urlencoded; charset=UTF-8",
            dataType: "json",
            success: function (result) {
                let appName = $("#app").find("option:selected").text();
                log.info(JSON.stringify(result));
                if (result.code == 200) {
                    let html = `<div style='padding:5%;'><p><b>产品：</b>${appName}</p><hr />`;
                    if (company == "torchsoft") {
                        html = `${html}<p><b>许可证数量：</b>${version}</p><hr />`;
                    } else {
                        html = `${html}<p><b>版本：</b>${version}</p><hr />`;
                    }
                    html = `${html}<p><b>key：</b>
                                <pre style='background: black;color:#66FF66;padding:5%;'>${result.data.key}</pre>
                            </p><hr /></div>`;

                    //自定页
                    layer.open({
                        // 在默认状态下，layer是宽高都自适应的，但当你只想定义宽度时，你可以area: '500px'，高度仍然是自适应的。
                        // 当你宽高都要定义时，你可以area: ['500px', '300px']
                        //area: [keyAreaWidthWidth],
                        type: 1,
                        icon: 1,
                        // 样式类名,目前layer内置的skin有：layui-layer-lan、layui-layer-molv
                        skin: 'layui-layer-lan',
                        // 关闭按钮
                        closeBtn: 1,
                        anim: 2,
                        // 点击遮罩是否关闭弹窗
                        shadeClose: false,
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

// =======================================  下载Netsarang  ======================================

window.netSarangDownload = function netSarangDownload() {
    let company = $("#company").val();
    if (company != "netsarang") {
        //提示层
        layer.msg("只提供NetSarang的产品最新版本下载", {icon: 5});
        return;
    }
    let app = $("#app").val();
    let version = $("#version").val();
    if (version != "latest") {
        //提示层
        layer.msg("只提供NetSarang的产品最新版本下载", {icon: 5});
        return;
    }
    //加载层,0代表加载的风格，支持0-2,0.5透明度的白色背景
    let index = layer.load(0, {shade: [0.5, '#fff']});
    $.ajax({
        url: "/getXshellUrl",
        type: "POST",
        data: {app: app, version: version},
        contentType: "application/x-www-form-urlencoded; charset=UTF-8",
        dataType: "json",
        success: function (result) {
            layer.close(index);
            let appName = $("#app").find("option:selected").text();
            if (result.code == 200) {
                let html = `<div style='padding:5%;text-align:center;word-wrap:break-word;'>
                                <p><b>${appName} 下载地址：</b></p>
                                <p><a href='${result.data.url}' target='_blank'>${result.data.url}</a></p>
                            </div>`;
                //自定页
                layer.open({
                    type: 1,
                    icon: 1,
                    // 样式类名,目前layer内置的skin有：layui-layer-lan、layui-layer-molv
                    skin: 'layui-layer-lan',
                    // 关闭按钮
                    closeBtn: 1,
                    anim: 2,
                    // 点击遮罩是否关闭弹窗
                    shadeClose: false,
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

// ====================================  格式化NGINX配置  ========================================


/**
 * 设置美化代码方式
 *
 * @return
 */
window.beautificationClick = function beautificationClick(event) {
    let value = $(event).val();
    if ("online" == value) {
        $("#indent-way").hide();
    } else if ("offline" == value) {
        $("#indent-way").show()
    }
}


/***
 * 设置缩进选中
 *
 * @return
 */
window.indentWayButton = function indentWayButton(event) {
    // 先去掉选中的
    $(".pure-button-active").removeClass('pure-button-active');
    //每次点击的时候，将当前的元素切换active样式
    $(event).addClass('pure-button-active');
}

/**
 * 点击美化按钮
 *
 * @return
 */
window.beautifyCode = function beautifyCode() {
    let beautification = $("input[name='beautification']:checked").val();
    let code = $("#text-code").val();
    if (util.isEmpty(code)) {
        layer.msg("请输入配置代码！");
        return;
    }
    if ("online" == beautification) {
        onlineBeautifier(code);
    } else if ("offline" == beautification) {
        let indentation = $(".pure-button-active").attr("id");
        if (util.isEmpty(indentation)) {
            layer.msg("请选择缩进方式！");
            return;
        }
        let indentCode = $("#indent-code").val();
        if ("space" == indentCode) {
            indentCode = "    ";
        } else if ("tab" == indentCode) {
            indentCode = "\t";
        }
        activateBeautifierListener(code, indentCode, indentation);
    }

}

/**
 * online美化Nginx配置
 *
 * @return
 */
function onlineBeautifier(code) {
    $.ajax({
        url: "/nginx-format-py",
        type: "POST",
        data: {code: code},
        contentType: "application/x-www-form-urlencoded; charset=UTF-8",
        dataType: "json",
        success: function (result) {

            if (result.code == 200) {
                beautifySuccess(result.data.contents);
            } else {
                //提示层
                layer.msg(result.msg, {icon: 5});
            }
        }
    })
}


/**
 * offline美化代码
 *
 * @param contents 配置代码
 * @param indentCode 缩进的代码
 * @param indentation 缩进方式
 *            indentWay1 按`server{\n`方式缩进(左花括号之后有一个空行)
 *            indentWay2 按`server{`方式缩进(左花括号之后无空行)
 *            indentWay3 按`server\n{`方式缩进（左花括号在新行中）
 * @return
 */
function activateBeautifierListener(contents, indentCode, indentation) {
    // 缩进代码
    INDENTATION = indentCode;
    // 缩进方式
    if (util.isEmpty(contents)) {
        layer.msg("请输入配置代码！");
        return;
    }
    modifyOptions({INDENTATION});
    // 将文件拆分为行，清理空格
    let cleanLines = clean_lines(contents);

    // 加入左括号（如果用户希望如此）默认为true
    let trailingBlankLines;
    if ("indentWay1" == indentation) {
        trailingBlankLines = true;
        modifyOptions({trailingBlankLines});
        cleanLines = join_opening_bracket(cleanLines);
    }
    // 加入左括号并且不要换行
    else if ("indentWay2" == indentation) {
        trailingBlankLines = false;
        modifyOptions({trailingBlankLines});
        cleanLines = join_opening_bracket(cleanLines);
    }
    // 执行最后的缩进
    cleanLines = array.trimSpace(perform_indentation(cleanLines, indentCode));

    // 将所有线条组合在一起
    let outputContents = cleanLines.join("\n");

    if ("indentWay2" == indentation) {
        outputContents = string.replaceBlank(outputContents);
    }
    // console.log(outputContents)
    // 将所有内容保存到文件中。
    // $("#text-code").val(outputContents);
    beautifySuccess(outputContents);
}


/**
 * 最后美化完成输出
 *
 * @return
 */
function beautifySuccess(contents) {
    let html = "<pre style='max-height: calc(93vh);background: black;color: #66FF66;margin: 0px;padding:" +
        " 10px;'>" + contents + "</pre>";
    //自定页
    layer.open({
        type: 1,
        icon: 1,
        skin: 'layui-layer-lan', //样式类名,目前layer内置的skin有：layui-layer-lan、layui-layer-molv
        closeBtn: 1, //关闭按钮
        anim: 2,
        shadeClose: true, //开启遮罩关闭
        title: false,
        content: html
    });
}