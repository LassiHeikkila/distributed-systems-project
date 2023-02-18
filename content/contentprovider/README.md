# contentprovider
`contentprovider` has copies of the video files and a database containing metadata. Database should be updated periodically, but it is not necessary to update it in real-time. Original data is controlled by `contentmanager`, which will distribute updated/new files and data when it wants.

## API
Need endpoints to:
- download file by id
- download subtitles
- get file metadata by id
- search files by metadata (title, length, type, producer, actors, etc.)
