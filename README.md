## drone-mella
`drone-mella` is a [drone](https://github.com/drone/drone) plugin for uploading files to [OwnCloud](https://owncloud.org/). It is baed on [mella](https://github.com/florianbeer/mella).

## Options

- `server`: OwnCloud server URL
- `remote_folder`: Folder in `server` where to store the file
- `files`: Array of files to be compressed.
- `owncloud_username`, `owncloud_password`: Credentials for `server`. You are advised to use environment variables `OWNCLOUD_USERNAME` and `OWNCLOUD_PASSWORD` in order to hide your credentials.

## Example configuration

```yml
pipeline:
  previousaction:
    [...]

  deploy:
    image: toroid/drone-mella
    server: https://owncloud.server.com
    remote_folder: "CREATE/THIS/BEFORE"
    files:
      - localFolder/*
```

## Notes

`drone-mella` will put `files` in a `.tgz` file named after the repository: `repoName_COMMIT-SHA.tgz` with `COMMIT-SHA` being the last 7 characters. If a tag hook is detected, drone will set `DRONE_TAG` and the file will be named: `repoName_TAG.tgz`.

Then this compressed file will be uploaded to `server/remote_folder`, provided the credentials are correct **and that the remote folder already exists**

## Contibuting

Don't hesitate to submit issues or pull requests.
