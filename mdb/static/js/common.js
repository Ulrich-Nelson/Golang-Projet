function formatDate(timestamping) {
    let date = new Date(timestamping * 1000);
    let months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    let year = date.getFullYear();
    let month = months[date.getMonth()];
    var day = date.getDate();
    return day + ' ' + month + ' ' + year;
}
