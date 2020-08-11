function generateUuidPart(e) {
    for (var t = "", r = 0; r < e; r++)
        t += parseInt(16 * Math.random()).toString(16).toUpperCase()
    return formatNum(t, e)
}

function generateUuid() {
    var e = generateUuidPart(8)
        , t = generateUuidPart(4)
        , r = generateUuidPart(4)
        , n = generateUuidPart(4)
        , o = generateUuidPart(12)
        , i = (new Date).getTime()
    return e + "-" + t + "-" + r + "-" + n + "-" + o + formatNum((i % 1e5).toString(), 5) + "infoc"
}

function formatNum(e, t) {
    var r = "";
    if (e.length < t)
        for (var n = 0; n < t - e.length; n++)
            r += "0"
    return r + e
}

function sendPV(e, t) {
    var r = 0 < arguments.length && void 0 !== e ? e : ""
        , n = 1 < arguments.length && void 0 !== t ? t : ""
    this.todo(r, n)

}

// console.log(generateUuid())

// console.log(sendPV());
