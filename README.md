## Intro

This project serves as the persistent logging backend to the lattice-prototype frontend.
It is built with Golang [Echo](https://echo.labstack.com/).
It can both persist every individual logging statement and batches of these statements.
Individual statements are stored to the root's local [logs.jsonl](logs.jsonl) file.
Batch statements are stored to the root's local [localStorageLogs.jsonl](localStorageLogs.jsonl) file.
Note that both files are not uploaded to git.

## Starting

Run the following command with setup Golang resources in your terminal of choice:

```shell
go build lattice-logging/cmd
```

### Starting the Server
Per default running above command will result in starting the persistent server feature.

### Starting the Analysis Tool
To run the analysis tool, you simply have to provide an environment variable
to the starting command. The variable must have the name `MODE` and accepts either 
`server` or `analyze` as proper values. The analysis tool is started by setting the `MODE`
to `analyze` accordingly.