# Kibana Object Provider

Provides a Logz.io kibana object resource. This can be used to export and import kibana objects.

* Learn more about kibana objects in the [Logz.io Docs](https://docs.logz.io/api/#tag/Import-or-export-Kibana-objects)

**Note:**
- The import operation is not available for this resource.
- This resource wraps the import/export API.
- The DELETE operation just deletes the resource from the state. To actually delete the object you'll need to manually delete it from the app.

## Example Usage

```hcl
variable "api_token" {
  type = "string"
  description = "your logzio API token"
}

provider "logzio" {
  api_token = var.api_token
}

resource "logzio_kibana_object" "my_search" {
  kibana_version = "7.2.1"
  data = file("/path/to/your/object/config.json")
}
```

### Example of config.json

```json
{
  "_index": "logzioCustomerIndex*",
  "_type": "_doc",
  "_id": "search:tf-provider-test-search",
  "_score": 1.290984,
  "_source": {
    "search": {
      "hits": 0,
      "columns": [
        "message"
      ],
      "description": "",
      "sort": [
        "@timestamp",
        "desc"
      ],
      "id": "tf-provider-test-search",
      "title": "tf provider test create search",
      "version": 1,
      "_updatedAt": 1561454443631,
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"highlight\":{\"pre_tags\":[\"@kibana-highlighted-field@\"],\"post_tags\":[\"@/kibana-highlighted-field@\"],\"fields\":{\"*\":{}},\"fragment_size\":2147483647},\"filter\":[],\"query\":{\"query\":\"type: tf-provider-test\",\"language\":\"lucene\"},\"source\":{\"excludes\":[]},\"highlightAll\":true,\"version\":true,\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
      },
      "panelsJSON": "[]"
    },
    "type": "search",
    "updated_at": 1561454443631,
    "references": [
      {
        "name": "kibanaSavedObjectMeta.searchSourceJSON.index",
        "type": "index-pattern",
        "id": "logzioCustomerIndex*"
      }
    ],
    "id": "tf-provider-test-search"
  }
}
```

## Argument Reference

* `kibana_version` - (String) The version of Kibana used at the time of export.
* `data` - (String) Exported Kibana objects. Should be a valid JSON that was retrieved from an export operation of the API.
