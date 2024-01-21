function redirectCreate() {
    var url = new URL(window.location.href)
    var categoryID = url.searchParams.get("category")

    var createNewsUrl = new URL(window.location.origin + "/create")
    if (categoryID) {
        createNewsUrl.searchParams.set("category", categoryID)
    }
    window.open(createNewsUrl, "_self")
}
