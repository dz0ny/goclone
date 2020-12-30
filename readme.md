# QuickCheck

Crawl domain and find missing assets(css, js, img with srcset) or links.

## Usage 

```
Usage of ./quickcheck:
  -allowed string
        Comma separated list of allowed domains to crawl
  -max int
        Max depth to crawl (default 100)
  -report string
        Path to the report file (default "/tmp/report.json")
  -url string
        URL to the website you want to check (default "https:/google.si")
```

## Example

```
./quickcheck -url=https://woocart.com -allowed=static.woocart.com,blogwoocartcom-e6da.kxcdn.com,fanstatic.niteo.co,blog.woocart.com,static.woocart.com/
{
 "https://static.woocart.com/images/10-support/solve_issues.png?v=1609318219": 1,
 "https://static.woocart.com/images/10-support/support?v=1609318219": 1,
 "https://static.woocart.com/images/content_images/blog/dashboard.png?v=1609318219": 1,
 "https://static.woocart.com/images/content_images/blog/themes_sidebar.png?v=1609318219": 1
}⏎
```