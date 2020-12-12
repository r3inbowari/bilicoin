function generateUuidPart(c) {
    for (var a = "", b = 0; b < c; b++) {
        a += parseInt(16 * Math.random()).toString(16).toUpperCase()
    }
    return formatNum(a, c)
}

function generateUuid() {
    var d = generateUuidPart(8), b = generateUuidPart(4), c = generateUuidPart(4), g = generateUuidPart(4),
        f = generateUuidPart(12), a = (new Date).getTime();
    return d + "-" + b + "-" + c + "-" + g + "-" + f + formatNum((a % 100000).toString(), 5) + "infoc"
}

function formatNum(c, a) {
    var b = "";
    if (c.length < a) {
        for (var d = 0; d < a - c.length; d++) {
            b += "0"
        }
    }
    return b + c
};
