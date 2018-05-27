package pdl

//go:generate go run gen.go -o har.go

// HAR is the PDL formatted definition of HAR types.
const HAR = `version
  major 1
  minor 3

# HTTP Archive Format
domain HAR

  # Contains info about a request coming from browser cache.
  type Cache extends object
    properties
      # State of a cache entry before the request. Leave out this field if the information is not available.
      optional CacheData beforeRequest
      # State of a cache entry after the request. Leave out this field if the information is not available.
      optional CacheData afterRequest
      # A comment provided by the user or the application.
      optional string comment

  # Describes the cache data for beforeRequest and afterRequest.
  type CacheData extends object
    properties
      # Expiration time of the cache entry.
      optional string expires
      # The last time the cache entry was opened.
      string lastAccess
      # Etag
      string eTag
      # The number of times the cache entry has been opened.
      integer hitCount
      # A comment provided by the user or the application.
      optional string comment

  # Describes details about response content (embedded in <response> object).
  type Content extends object
    properties
      # Length of the returned content in bytes. Should be equal to response.bodySize if there is no compression and bigger when the content has been compressed.
      integer size
      # Number of bytes saved. Leave out this field if the information is not available.
      optional integer compression
      # MIME type of the response text (value of the Content-Type response header). The charset attribute of the MIME type is included (if available).
      string mimeType
      # Response body sent from the server or loaded from the browser cache. This field is populated with textual content only. The text field is either HTTP decoded text or a encoded (e.g. "base64") representation of the response body. Leave out this field if the information is not available.
      optional string text
      # Encoding used for response text field e.g "base64". Leave out this field if the text field is HTTP decoded (decompressed & unchunked), than trans-coded from its original character set into UTF-8.
      optional string encoding
      # A comment provided by the user or the application.
      optional string comment

  # Contains list of all cookies (used in <request> and <response> objects).
  type Cookie extends object
    properties
      # The name of the cookie.
      string name
      # The cookie value.
      string value
      # The path pertaining to the cookie.
      optional string path
      # The host of the cookie.
      optional string domain
      # Cookie expiration time. (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD, e.g. 2009-07-24T19:20:30.123+02:00).
      optional string expires
      # Set to true if the cookie is HTTP only, false otherwise.
      optional boolean httpOnly
      # True if the cookie was transmitted over ssl, false otherwise.
      optional boolean secure
      # A comment provided by the user or the application.
      optional string comment

  # Creator and browser objects share the same structure.
  type Creator extends object
    properties
      # Name of the application/browser used to export the log.
      string name
      # Version of the application/browser used to export the log.
      string version
      # A comment provided by the user or the application.
      optional string comment

  # Represents an array with all exported HTTP requests. Sorting entries by startedDateTime (starting from the oldest) is preferred way how to export data since it can make importing faster. However the reader application should always make sure the array is sorted (if required for the import).
  type Entry extends object
    properties
      # Reference to the parent page. Leave out this field if the application does not support grouping by pages.
      optional string pageref
      # Date and time stamp of the request start (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD).
      string startedDateTime
      # Total elapsed time of the request in milliseconds. This is the sum of all timings available in the timings object (i.e. not including -1 values) .
      number time
      # Detailed info about the request.
      Request request
      # Detailed info about the response.
      Response response
      # Info about cache usage.
      Cache cache
      # Detailed timing info about request/response round trip.
      Timings timings
      # IP address of the server that was connected (result of DNS resolution).
      optional string serverIPAddress
      # Unique ID of the parent TCP/IP connection, can be the client or server port number. Note that a port number doesn't have to be unique identifier in cases where the port is shared for more connections. If the port isn't available for the application, any other unique connection ID can be used instead (e.g. connection index). Leave out this field if the application doesn't support this info.
      optional string connection
      # A comment provided by the user or the application.
      optional string comment

  # Parent container for HAR log.
  type HAR extends object
    properties
      Log log

  # Represents the root of exported data.
  type Log extends object
    properties
      # Version number of the format. If empty, string "1.1" is assumed by default.
      string version
      # Name and version info of the log creator application.
      Creator creator
      # Name and version info of used browser.
      optional Creator browser
      # List of all exported (tracked) pages. Leave out this field if the application does not support grouping by pages.
      optional array of Page pages
      # List of all exported (tracked) requests.
      array of Entry entries
      # A comment provided by the user or the application.
      optional string comment

  # Describes a name/value pair.
  type NameValuePair extends object
    properties
      # Name of the pair.
      string name
      # Value of the pair.
      string value
      # A comment provided by the user or the application.
      optional string comment

  # Represents list of exported pages.
  type Page extends object
    properties
      # Date and time stamp for the beginning of the page load (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD, e.g. 2009-07-24T19:20:30.45+01:00).
      string startedDateTime
      # Unique identifier of a page within the <log>. Entries use it to refer the parent page.
      string id
      # Page title.
      string title
      # Detailed timing info about page load.
      PageTimings pageTimings
      # A comment provided by the user or the application.
      optional string comment

  # Describes timings for various events (states) fired during the page load. All times are specified in milliseconds. If a time info is not available appropriate field is set to -1.
  type PageTimings extends object
    properties
      # Content of the page loaded. Number of milliseconds since page load started (page.startedDateTime). Use -1 if the timing does not apply to the current request.
      optional number onContentLoad
      # Page is loaded (onLoad event fired). Number of milliseconds since page load started (page.startedDateTime). Use -1 if the timing does not apply to the current request.
      optional number onLoad
      # A comment provided by the user or the application.
      optional string comment

  # List of posted parameters, if any (embedded in <postData> object).
  type Param extends object
    properties
      # name of a posted parameter.
      string name
      # value of a posted parameter or content of a posted file.
      optional string value
      # name of a posted file.
      optional string fileName
      # content type of a posted file.
      optional string contentType
      # A comment provided by the user or the application.
      optional string comment

  # Describes posted data, if any (embedded in <request> object).
  type PostData extends object
    properties
      # Mime type of posted data.
      string mimeType
      # List of posted parameters (in case of URL encoded parameters).
      array of Param params
      # Plain text posted data
      string text
      # A comment provided by the user or the application.
      optional string comment

  # Contains detailed info about performed request.
  type Request extends object
    properties
      # Request method (GET, POST, ...).
      string method
      # Absolute URL of the request (fragments are not included).
      string url
      # Request HTTP Version.
      string httpVersion
      # List of cookie objects.
      array of Cookie cookies
      # List of header objects.
      array of NameValuePair headers
      # List of query parameter objects.
      array of NameValuePair queryString
      # Posted data info.
      optional PostData postData
      # Total number of bytes from the start of the HTTP request message until (and including) the double CRLF before the body. Set to -1 if the info is not available.
      integer headersSize
      # Size of the request body (POST data payload) in bytes. Set to -1 if the info is not available.
      integer bodySize
      # A comment provided by the user or the application.
      optional string comment

  # Contains detailed info about the response.
  type Response extends object
    properties
      # Response status.
      integer status
      # Response status description.
      string statusText
      # Response HTTP Version.
      string httpVersion
      # List of cookie objects.
      array of Cookie cookies
      # List of header objects.
      array of NameValuePair headers
      # Details about the response body.
      Content content
      # Redirection target URL from the Location response header.
      string redirectURL
      # Total number of bytes from the start of the HTTP response message until (and including) the double CRLF before the body. Set to -1 if the info is not available.
      integer headersSize
      # Size of the received response body in bytes. Set to zero in case of responses coming from the cache (304). Set to -1 if the info is not available.
      integer bodySize
      # A comment provided by the user or the application.
      optional string comment

  # Describes various phases within request-response round trip. All times are specified in milliseconds.
  type Timings extends object
    properties
      # Time spent in a queue waiting for a network connection. Use -1 if the timing does not apply to the current request.
      optional number blocked
      # DNS resolution time. The time required to resolve a host name. Use -1 if the timing does not apply to the current request.
      optional number dns
      # Time required to create TCP connection. Use -1 if the timing does not apply to the current request.
      optional number connect
      # Time required to send HTTP request to the server.
      number send
      # Waiting for a response from the server.
      number wait
      # Time required to read entire response from the server (or cache).
      number receive
      # Time required for SSL/TLS negotiation. If this field is defined then the time is also included in the connect field (to ensure backward compatibility with HAR 1.1). Use -1 if the timing does not apply to the current request.
      optional number ssl
      # A comment provided by the user or the application.
      optional string comment
`
