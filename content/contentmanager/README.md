# content-manager
This service processes files uploaded by admin.
It automatically creates transcoding and downscaling jobs when new videos are uploaded.

Transcoding/downscaling job should not be made for files which have a non-null `originalFileID` value.

`content-manager` holds the primary copies of videos and metadata about them.
It distributes copies to instances of `content-provider`.
Distribution of copies does not need to be instantaneous.

Distribution mechanism could be as simple as `rsync`.
