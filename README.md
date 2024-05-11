# Sky Lantern

## Overview

Sky Lantern is an HTTP/HTTPS based multi-source downloader that can download a file from multiple sources simultaneously.
This project simulates a scenario where a file is split into multiple chunks and distributed across different servers (e.g. BitTorrent).

**Since we don't have tracker servers in this project, the user must provide the URLs of all the chunk URLs to simulate a manifest file.**
(This project also not support `Transfer-Encoding: chunked` response as our chunks are on different servers)

## Usage

To download a file, run the following command:

```bash
# Support arguments:
# -keep-chunks: Keep the downloaded chunks as cache, will not download the same chunk again.
# -output: The output file path
# -debug: Print debug information

# Example:
go run main.go -keep-chunks -output result.txt -debug <chunk1_url> <chunk2_url> <chunk3_url> ...
# or use the provided run.sh script with sample URLs
bash ./run.sh
```

## Project Structure

- `main.go`: The entry point of the application, handling user input and starting the download process.
- `downloader/downloader.go`: Contains the main logic for the multi-source downloader, including chunk reassembly and hash verification.
- `utils/utils.go`: Provides utility functions, such as making HTTP requests and calculating file hashes.

## Change Log

- **2024-05-11**: My initial thought was we can grab the chunks URL from the response header.
But that will only work for custom servers. So, I decided to let the user provide the URLs of all the chunks to make it more generic.

## Future Improvements

- In current implementation, we first grab the chunk headers to get the file size and the ETag(hash), it seems not necessary for the project.
  But this is for a real-world scenario where download behavior might depend on the information in the headers. (Kind of like `stat` operation)
- When we write our chunks to the file, we can leverage cursor to write the chunks to the correct position.
- Use channels to handle the download of chunks concurrently. Currently, we are using Lock-Free shared memory to store the chunks.
- Implement a tracker server to manage the distribution of chunks.
