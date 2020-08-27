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
 * 请求方式（OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, PATCH）
 * https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Methods
 *
 * @type {{TRACE: string, HEAD: string, DELETE: string, POST: string, GET: string, PATCH: string, OPTIONS: string, PUT: string}}
 */
const METHOD = {
    GET: "GET",
    HEAD: "HEAD",
    POST: "POST",
    PUT: "PUT",
    DELETE: "DELETE",
    CONNECT: "CONNECT",
    OPTIONS: "OPTIONS",
    TRACE: "TRACE",
    PATCH: "PATCH",
}

/**
 * Content-Type请求数据类型，告诉接收方，我发什么类型的数据。
 * https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Type
 * https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types
 * https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
 *
 * application/x-www-form-urlencoded：数据被编码为名称/值对。这是标准的编码格式。默认使用此类型。
 * multipart/form-data：数据被编码为一条消息，页上的每个控件对应消息中的一个部分。
 * text/plain：数据以纯文本形式(text/json/xml/html)进行编码，其中不含任何控件或格式字符。postman软件里标的是RAW。
 *
 * @type {{FORM_DATA: string, URLENCODED: string, TEXT_PLAIN: string}}
 */
const CONTENT_TYPE = {
    URLENCODED: "application/x-www-form-urlencoded", FORM_DATA: "multipart/form-data", TEXT_PLAIN: "text/plain"
}

/**
 * XMLHttpRequest预期服务器返回数据类型，并根据此值进行本地解析
 * https://developer.mozilla.org/zh-CN/docs/Web/API/XMLHttpRequest/responseType
 *
 * @type {{ARRAY_BUFFER: string, BLOB: string, MS_STREAM: string, DOCUMENT: string, TEXT: string, JSON: string}}
 */
const RESPONSE_TYPE = {
    TEXT: "text", ARRAY_BUFFER: "arraybuffer", BLOB: "blob", DOCUMENT: "document", JSON: "json", MS_STREAM: "ms-stream"
}


/**
 * js封装ajax请求 https://developer.mozilla.org/zh-CN/docs/Web/API/XMLHttpRequest
 * 使用new XMLHttpRequest 创建请求对象,所以不考虑低端IE浏览器(IE6及以下不支持XMLHttpRequest)
 * 注意:请求参数如果包含日期类型.是否能请求成功需要后台接口配合
 *
 *   url：       请求路径
 *   method：    请求方式（OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, PATCH）
 *   data：      是作为请求主体被发送的数据,只适用于这些请求方法 'PUT','POST','PATCH'
 *   contentType：  请求数据类型(application/x-www-form-urlencoded,multipart/form-data,text/plain)
 *   responseType： 响应的数据类型（text,arraybuffer,blob,document,json,ms-stream）
 *   timeout：      超时时间，0表示不设置超时
 *
 * @param config 请求参数
 */
const ajax = (settings = {}) => {
    // 初始化请求参数
    let config = Object.assign({
        url: '',
        method: settings.type || settings.method || METHOD.GET,
        // string 期望的返回数据类型:'json' 'text' 'document' ...
        responseType: settings.dataType || settings.responseType || RESPONSE_TYPE.JSON,
        async: true, //  boolean true:异步请求 false:同步请求 required
        data: null, // any 请求参数,data需要和请求头Content-Type对应
        headers: {},
        timeout: settings.timeout || 1000, // 超时时间:0表示不设置超时
        beforeSend: (xhr) => {

        },
        success: (result, status, xhr) => {

        },
        error: (xhr, status, error) => {

        },
        complete: (xhr, status) => {

        }
    }, settings);

    if (!config.headers["Content-Type"]) {
        // 服务器会根据此值解析参数，同时在返回时也指定此值
        config.headers["Content-Type"] = settings.contentType || config.headers["Content-Type"]
            || config.headers["content-type"] || CONTENT_TYPE.URLENCODED;
    }
    // 参数验证
    if (!config.url) {
        throw new TypeError("ajax请求：url参数不正确");
    }
    if (!config.method) {
        throw new TypeError("ajax请求：type或method参数不正确");
    }
    if (!config.responseType) {
        throw new TypeError("ajax请求：dataType或responseType参数不正确");
    }
    if (!config.headers || !config.headers["Content-Type"]) {
        throw new TypeError("ajax请求：Content-Type参数不正确");
    }
    // 创建XMLHttpRequest请求对象
    let xhr = new XMLHttpRequest();

    // 请求开始回调函数，对应xhr.loadstart
    xhr.addEventListener('loadstart', e => {
        config.beforeSend(xhr, e);
    });
    // 请求成功回调函数，对应xhr.onload
    xhr.addEventListener('load', e => {
        if ((xhr.status >= 200 && xhr.status < 300) || xhr.status === 304) {
            if (xhr.responseType === 'text') {
                config.success(xhr.responseText, xhr.status, xhr);
            } else if (xhr.responseType === 'document') {
                config.success(xhr.responseXML, xhr.status, xhr);
            } else {
                config.success(xhr.response, xhr.status, xhr);
            }
        } else {
            config.error(xhr, xhr.status, e);
        }
    });
    // 请求结束，对应xhr.onloadend
    xhr.addEventListener('loadend', e => {
        config.complete(xhr, xhr.status, e);
    });
    // 请求出错，对应xhr.onerror
    xhr.addEventListener('error', e => {
        config.error(xhr, xhr.status, e);
    });
    // 请求超时，对应xhr.ontimeout
    xhr.addEventListener('timeout', e => {
        config.error(xhr, 408, e);
    });

    let method = config.method.toUpperCase();
    // 如果是"简单"请求,则把data参数组装在url上
    if ((method === 'GET' || method === 'DELETE') && config.data) {
        let paramsStr;
        if ((config.data).constructor === Object) {
            let paramsArr = [];
            Object.keys(config.data).forEach(key => {
                paramsArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(config.data[key])}`);
            });
            paramsStr = paramsArr.join('&');
        } else if ((config.data).constructor === String) {
            paramsStr = config.data;
        } else if ((config.data).constructor === Array) {
            paramsStr = config.data.join("&");
        } else {
            throw new TypeError("ajax请求：数据类型错误！");
        }
        config.url += (config.url.indexOf('?') !== -1) ? paramsStr : '?' + paramsStr;
    }

    // 初始化请求
    xhr.open(config.method, config.url, config.async);
    // 设置请求头，必须要放到open()后面
    for (const key of Object.keys(config.headers)) {
        xhr.setRequestHeader(key, config.headers[key]);
    }
    // 设置超时时间
    if (config.timeout) {
        xhr.timeout = config.timeout;
    }
    // 设置预期服务器返回数据类型，并进行本地解析
    xhr.responseType = config.responseType;

    // 请求参数类型需要和请求头Content-Type对应'PUT','POST','PATCH'
    if ((method === 'PUT' || method === 'POST' || method === 'PATCH') && config.data) {
        let ct = config.headers["Content-Type"].split(";")[0].toLocaleLowerCase();
        if (ct == "application/x-www-form-urlencoded") {
            if ((config.data).constructor !== Object) {
                throw new TypeError("ajax请求：application/x-www-form-urlencoded数据类型错误！");
            }
            let paramsArr = [];
            Object.keys(config.data).forEach(key => {
                paramsArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(config.data[key])}`);
            });
            xhr.send(paramsArr.join('&'));
        } else if (ct == "multipart/form-data") {
            if ((config.data).constructor !== Object) {
                throw new TypeError("ajax请求：multipart/form-data数据类型错误！");
            }
            let formData = new FormData();
            Object.keys(config.data).forEach(key => {
                formData.append(key, config.data[key]);
            });
            xhr.send(formData);
        } else if (ct == "text/plain") {
            if ((config.data).constructor !== String) {
                throw new TypeError("ajax请求：text/plain数据类型错误！");
            }
            xhr.send(config.data);
        } else if (ct == "application/json") {
            if ((config.data).constructor === String) {
                try {
                    JSON.parse(config.data);
                    xhr.send(config.data);
                } catch (e) {
                    throw new TypeError("ajax请求：application/json数据类型错误！");
                }
            } else if ((config.data).constructor === Array || (config.data).constructor === Object) {
                xhr.send(JSON.stringify(config.data));
            } else {
                throw new TypeError("ajax请求：application/json数据类型错误！");
            }
        } else {
            throw new TypeError("ajax请求：数据类型错误！");
        }
    } else {
        xhr.send();
    }
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
        ajax({
            url: url,
            method: METHOD.POST,
            data: params,
            responseType: RESPONSE_TYPE.BLOB,
            success: (result, status, xhr) => {
                // console.log(xhr.getAllResponseHeaders())
                //从response的headers中获取filename, 后端response.setHeader("Content-Disposition", "attachment; filename=xxxx.xxx") 设置的文件名;
                let contentDisposition = xhr.getResponseHeader('Content-Disposition') || xhr.getResponseHeader('content-disposition');
                let contentType = xhr.getResponseHeader('Content-Type') || xhr.getResponseHeader('content-type') || 'application/octet-stream';
                // let contentLength = result.headers["Content-Length"] || result.headers["content-length"];
                let filename;
                // 如果从Content-Disposition中取到的文件名不为空
                if (!util.isEmpty(contentDisposition)) {
                    // 取出文件名，这里正则注意顺序 (.*)在|前如果有;号那么永远都会是真 把分号以及后面的字符取到
                    let reg = new RegExp("(?<=filename=)((.*)(?=;|%3B)|(.*))").exec(contentDisposition);
                    // 取文件名信息中的文件名,替换掉文件名中多余的符号
                    filename = reg[1].replaceAll("\\\\|/|\"", "");
                } else {
                    let urls = url.split("/");
                    filename = urls[urls.length - 1];
                }
                // 解决中文乱码，编码格式
                filename = decodeURI(escape(filename));
                let downloadElement = document.createElement('a');
                downloadElement.style.display = 'none';
                // 创建下载的链接
                downloadElement.href = URL.createObjectURL(new Blob([xhr.response], {type: contentType}));
                // 下载后文件名
                downloadElement.download = filename;
                // 点击下载
                downloadElement.click();
                // 释放掉blob对象
                URL.revokeObjectURL(downloadElement.href);
            },
            error: (xhr, status, error) => {
                // 如果服务器自定义错误返回
                if (xhr.response && xhr.response.type === 'application/json') {
                    let reader = new FileReader();
                    reader.readAsText(xhr.response, 'utf-8');
                    reader.onload = () => {
                        console.log(reader.result)
                        reject(JSON.parse(reader.result).message);
                    }
                    return;
                }
                reject(xhr.response);
            }
        })
    })
}


/**
 * export default 服从 ES6 的规范,补充：default 其实是别名
 * module.exports 服从 CommonJS 规范 https://javascript.ruanyifeng.com/nodejs/module.html
 * 一般导出一个属性或者对象用 export default
 * 一般导出模块或者说文件使用 module.exports
 *
 * import from 服从ES6规范,在编译器生效
 * require 服从ES5 规范，在运行期生效
 * 目前 vue 编译都是依赖label 插件，最终都转化为ES5
 *
 * @return 将方法、变量暴露出去
 */
export default {
    METHOD,
    CONTENT_TYPE,
    RESPONSE_TYPE,
    ajax,
    download
}