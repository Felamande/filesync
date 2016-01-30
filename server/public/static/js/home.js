$(function () {
    new Vue({
        el: ".container-fluid",
        data: {
            msg: "one"
        },
        components: {
            "affix": VueStrap.affix,
            "alert": VueStrap.alert,
            "tab": VueStrap.tab,
            "tabs": VueStrap.tabset
        }
    })
})