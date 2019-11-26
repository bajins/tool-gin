/**
 * @Description:
 * @Author: bajins www.bajins.com
 * @File: http.js
 * @Version: 1.0.0
 * @Time: 2019/9/12 11:29
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */

import util from "./util.js";

/**
 * 请求方式（OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, CONNECT）
 *
 * @type {{TRACE: string, HEAD: string, DELETE: string, POST: string, GET: string, CONNECT: string, OPTIONS: string, PUT: string}}
 */
const METHOD = {
    OPTIONS: "OPTIONS",
    GET: "GET",
    HEAD: "HEAD",
    POST: "POST",
    PUT: "PUT",
    DELETE: "DELETE",
    TRACE: "TRACE",
    CONNECT: "CONNECT"
}

/**
 * 请求数据类型,告诉服务器，我要发什么类型的数据。
 *
 * application/x-www-form-urlencoded：数据被编码为名称/值对。这是标准的编码格式。默认使用此类型。
 * multipart/form-data：数据被编码为一条消息，页上的每个控件对应消息中的一个部分。
 * text/plain：数据以纯文本形式(text/json/xml/html)进行编码，其中不含任何控件或格式字符。postman软件里标的是RAW。
 *
 * @type {{FORM_DATA: string, URLENCODED: string, TEXT_PLAIN: string}}
 */
const CONTENT_TYPE = {
    URLENCODED: "application/x-www-form-urlencoded",
    FORM_DATA: "multipart/form-data",
    TEXT_PLAIN: "text/plain"
}

/**
 * 预期服务器返回的数据类型（对应请求头中的Accept），如果是下载文件则指定RESPONSE_TYPE
 *
 * 如果没有指定，那么会自动推断是返回 XML，还是JSON，还是script，还是String。
 * xml: 返回 XML 文档。
 * html: 返回纯文本 HTML 信息；包含的 script 标签会在插入 dom 时执行。
 * script: 返回纯文本 JavaScript 代码。不会自动缓存结果。除非设置了 “cache” 参数。
 * 注意：在远程请求时(不在同一个域下)，所有 POST 请求都将转为 GET 请求。（因为将使用 DOM 的 script标签来加载）
 * json: 返回 JSON 数据 。
 * jsonp: JSONP 格式。使用 JSONP 形式调用函数时，如 “myurl?callback=?” jQuery 将自动替换 ? 为正确的函数名，以执行回调函数。
 * text: 返回纯文本字符串
 *
 * @type {{SCRIPTY: string, JSONP: string, XML: string, JSON: string, TEXT: string, HTML: string}}
 */
const DATA_TYPE = {
    JSON: "json", TEXT: "text", XML: "xml", HTML: "html", SCRIPTY: "script", JSONP: "jsonp"
}

/**
 * 响应的数据类型
 *
 *   ""    将 responseType 设为空字符串与设置为"text"相同， 是默认类型 （实际上是 DOMString）。
 *  "arraybuffer"    response 是一个包含二进制数据的 JavaScript ArrayBuffer 。
 *  "blob"    response 是一个包含二进制数据的 Blob 对象 。
 *  "document"    response 是一个 HTML Document 或 XML XMLDocument ，这取决于接收到的数据的 MIME 类型。
 *  "json"    response 是一个 JavaScript 对象。这个对象是通过将接收到的数据类型视为 JSON 解析得到的。
 *  "text"    response 是包含在 DOMString 对象中的文本。
 *  "moz-chunked-arraybuffer" 与"arraybuffer"相似，但是数据会被接收到一个流中。
 *         使用此响应类型时，响应中的值仅在 progress 事件的处理程序中可用，并且只包含上一次响应 progress 事件以后收到的数据，
 *         而不是自请求发送以来收到的所有数据。在 progress 事件处理时访问 response 将返回到目前为止收到的数据。
 *         在 progress 事件处理程序之外访问， response 的值会始终为 null 。
 *  "ms-stream"  response 是下载流的一部分；此响应类型仅允许下载请求，并且仅受Internet Explorer支持。
 *
 * @type {{ARRAYBUFFER: string, BLOB: string, MS_STREAM: string, DOCUMENT: string, TEXT: string, JSON: string}}
 */
const RESPONSE_TYPE = {
    TEXT: "text", ARRAY_BUFFER: "arraybuffer", BLOB: "blob", DOCUMENT: "document", JSON: "json", MS_STREAM: "ms-stream"
}


/**
 * 封装axios HTTP请求API为`Promise`方式
 * 使用方法：http.axiosRequest({obj对象的数据},url字符串：如果obj.url为空就取这里的值)
 *
 * @param url 请求路径
 * @param obj 有以下参数：
 *   url： 请求路径：如果obj.url为空就取这里的值
 *   method： 请求方式（OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, CONNECT）
 *   data： 是作为请求主体被发送的数据,只适用于这些请求方法 'PUT', 'POST', 和 'PATCH'
 *   params: 是即将与请求一起发送的 URL 参数
 *   contentType:  请求数据类型(application/x-www-form-urlencoded,multipart/form-data,text/plain)
 *   dataType： 返回数据类型（json,text,xml,html,script,jsonp）
 *   responseType： 响应的数据类型（text，arraybuffer,blob,document,json,ms-stream）
 *
 * @param url
 * @param obj
 * @returns {Promise<unknown>}
 */
const request = (url, obj) => {
    return new Promise((resolve, reject) => {
        $.ajax({
            url: obj.url || url,
            method: obj.method || METHOD.GET,
            data: obj.data || {},// 是作为请求主体被发送的数据,只适用于这些请求方法 'PUT', 'POST', 和 'PATCH'
            params: obj.data || {},// 是即将与请求一起发送的 URL 参数
            header: {
                'Content-Type': obj.contentType || CONTENT_TYPE.URLENCODED
            },
            dataType: obj.dataType || DATA_TYPE.JSON,
            responseType: obj.responseType || RESPONSE_TYPE.JSON,
        }).then(response => {
            resolve(response);
        }).catch((error) => {
            reject(error)
        })
    })
}


/**
 * 文件下载api封装
 *
 * @param url
 * @param params
 * @returns {Promise<unknown>}
 */
const download = (url, params) => {
    return new Promise((resolve, reject) => {
        request(url, {
            method: METHOD.POST,
            data: params,
            contentType: CONTENT_TYPE.URLENCODED,
            responseType: RESPONSE_TYPE.BLOB
        }).then(function (result, status, xhr) {

            // console.log(xhr.getAllResponseHeaders());
            // xhr.getResponseHeader('Content-Disposition');
            // 从response的headers中获取filename,
            // 后端response.setHeader("Content-Disposition", "attachment; filename=文件名");
            let contentDisposition = xhr.headers['Content-Disposition'];
            let filename = "";
            // 如果从Content-Disposition中取到的文件名不为空
            if (!util.isEmpty(contentDisposition)) {
                let reg = new RegExp("filename=([^;]+\\.[^\\.;]+);*");
                filename = reg.exec(contentDisposition)[1];
                // 取文件名信息中的文件名,替换掉文件名中多余的符号
                filename = filename.replaceAll("\\\\|/|\"", "");
            }
            let downloadElement = document.createElement('a');

            //这里res.data是返回的blob对象
            let blob = new Blob([result], {type: 'application/octet-stream;charset=utf-8'});
            // 创建下载的链接
            let href = window.URL.createObjectURL(blob);

            downloadElement.style.display = 'none';
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

        }).catch(function (err) {
            reject(err);
        })
    })
}


/**
 * export default 服从 ES6 的规范,补充：default 其实是别名
 * module.exports 服从CommonJS 规范
 * 一般导出一个属性或者对象用 export default
 * 一般导出模块或者说文件使用 module.exports
 *
 * import from 服从ES6规范,在编译器生效
 * require 服从ES5 规范，在运行期生效
 * 目前 vue 编译都是依赖label 插件，最终都转化为ES5
 *
 * @return 将方法、变量暴露出去
 * @Description
 * @author claer woytu.com
 * @date 2019/4/29 11:58
 */
export default {
    METHOD,
    CONTENT_TYPE,
    DATA_TYPE,
    RESPONSE_TYPE,
    request,
    download
}