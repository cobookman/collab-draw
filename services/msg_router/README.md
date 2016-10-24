### Deploying

Publishing the `subscribe` function in index.js

`gcloud alpha functions deploy subscribe --bucket [YOUR_BUCKET_NAME] --trigger-topic [YOUR_TOPIC_NAME]`


Getting the logs for cloud function:
`gcloud alpha functions get-logs subscribe`

