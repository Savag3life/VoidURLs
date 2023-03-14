### VoidURLs 
Simple URL shortener API written in GoLang with MongoDB as the backend. Resulting shortened URLs are simply `www.yourdomain.com/` followed by the API response from `/void/create`

#### Usage
Java & Golang Examples Provided Below
```java
OkHttpClient client = new OkHttpClient();

String url = "http://localhost:8080/void/create";
String json = "{\"url\":\"www.google.com\"}";

MediaType mediaType = MediaType.parse("application/json");
RequestBody requestBody = RequestBody.create(json, mediaType);

Request request = new Request.Builder()
        .url(url)
        .post(requestBody)
        .build();

try (Response response = client.newCall(request).execute()) {
    if (!response.isSuccessful()) {
        throw new IOException("Unexpected code " + response);
    }

    System.out.println(response.body().string());
} catch (IOException e) {
    e.printStackTrace();
}
```
```go
type RequestURL struct {
    InputURL string `json:"url"`
}

requestBody, err := json.Marshal(RequestURL{InputURL: "www.google.com"})
resp, err := http.Post("http://localhost:8080/void/create", "applicationjson", bytes.NewBuffer(requestBody))

defer resp.Body.Close()

var responseID string
if err := json.NewDecoder(resp.Body).Decode(&responseID); err != nil {
    fmt.Println("Error decoding response body:", err)
    return
}

fmt.Println("Short URL ID:", responseID)
```
#### Run
To run this project, 
`git clone https://github.com/Savag3life/VoidURLs.git`
`cd VoidURls`
`go run main.go`