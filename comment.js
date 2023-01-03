var apikey = 'AIzaSyASFbBIOoXju0Oz_xfprimUWUAx4FuiogI';
import fetch from 'node-fetch';
var list = 'PLhu18ozRJ5d3XLfoeUw6WQxGJAydvC-iL'
import fs from 'fs';

var result = []

async function GetId() {
    try {
        let res = await fetch(`https://www.googleapis.com/youtube/v3/playlistItems?key=${apikey}&playlistId=${list}&part=snippet&maxResults=50`)
        let json = await res.json();
        let items = await json.items
        let num = await items.length

        for (let i = 0; i < num; i++) {
            //動画ID、タイトル、サムネイルを取得する
            let Id = await items[i].snippet.resourceId.videoId
            let VideoTitle = await items[i].snippet.title
            let Image = await items[i].snippet.thumbnails.default
            await GetTimeStamp(
                {
                    "Id": Id,
                    "Title": VideoTitle,
                    "Image": Image
                }
            )
        }
    }
    catch (e) {
        console.log(e)
    }

    fs.writeFile('data.json', JSON.stringify(result, null, '    '), (err) => {
        if (err) console.log(`error!::${err}`);
    });
}

async function GetTimeStamp(props) {
    try {

        let res = await fetch(`https://www.googleapis.com/youtube/v3/commentThreads?key=${apikey}&part=snippet&videoId=${props.Id}`)
        let json = await res.json();
        let num = await json.items.length
        let url = await `https://www.youtube.com/watch?v=${props.Id}`

        for (let i = 0; i < num; i++) {
            let comment = await json.items[i].snippet.topLevelComment.snippet.textOriginal
            //コメントからタイムスタンプ部分を抜き出す
            let body = await (String(comment).match(/[0-9]{1,}:[0-9]{1,}:[0-9]{1,}(.*)(.*)|[0-9]{1,}:[0-9]{1,}(.*)(.*)/gi));
            let time = await (String(comment).match(/[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}/gi))
            let title = await (String(body).replace(/[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}/gi, ""));
            let titledata = await title.split(',');
            if (time.length > 5) {
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
                        'id': props.Id,
                        'videotitle': props.Title,
                        'image': props.Image,
                        'url': url
                    }
                })
                break;
            }
        }


    }
    catch {
        console.log('Failed.')
    }

}
GetId();

