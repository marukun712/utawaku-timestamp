var apikey = 'AIzaSyASFbBIOoXju0Oz_xfprimUWUAx4FuiogI';
import fetch from 'node-fetch';
var list = 'PLhu18ozRJ5d3XLfoeUw6WQxGJAydvC-iL'
import fs from 'fs';

var result = []

async function getid() {
    let res = await fetch(`https://www.googleapis.com/youtube/v3/playlistItems?key=${apikey}&playlistId=${list}&part=snippet&maxResults=50`)
    let json = await res.json();
    let items = await json.items
    let num = await items.length
    for (let i = 0; i < num; i++) {
        var id = await items[i].snippet.resourceId.videoId
        var videotitle = await items[i].snippet.title
        var image = await items[i].snippet.thumbnails.default
        var url = await `https://www.youtube.com/watch?v=${id}`
        let res = await fetch(`https://www.googleapis.com/youtube/v3/commentThreads?key=${apikey}&part=snippet&videoId=${id}`)
        let json = await res.json();
        let num = json.items?.length
        for (let i = 0; i < num; i++) {
            var comment = await json.items[i].snippet.topLevelComment.snippet.textOriginal
            var str = /[0-9]{1,}:[0-9]{1,}:[0-9]{1,}(.*)(.*)|[0-9]{1,}:[0-9]{1,}(.*)(.*)/gi
            let body = await (String(comment).match(str));
            let time = await (String(comment).match(/[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}/gi))
            let title = await (String(body).replace(/[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}/gi, ""));
            let titledata = title.split(',');
            if (time?.length < 5) {
            } else if (body === null) {
            } else {
                let stampdata = []
                for (let i = 0; i < time.length; i++) {
                    stampdata.push({
                        item: {
                            'time': time[i],
                            'title': titledata[i]
                        }
                    })
                }
                result.push({
                    item: {
                        'timestamp': stampdata,
                        'id': id,
                        'videotitle': videotitle,
                        'image': image,
                        'url': url
                    }
                })
                break;
            }

        }


    }
    fs.writeFile('data.json', JSON.stringify(result, null, '    '), (err) => {
        if (err) console.log(`error!::${err}`);
    });
}
getid();

