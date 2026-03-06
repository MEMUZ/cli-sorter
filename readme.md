# 📂 CLI Sorter

A simple and fast CLI tool written in Go that automatically organizes files in a directory by their file type.
It scans a folder and moves files into categorized directories like:

- images
- videos
- documents
- archives
- other

Perfect for quickly cleaning up messy Downloads folders.

## ✨ Features

- ⚡ Fast file sorting
- 📁 Automatic folder creation
- 🔍 Dry-run mode to preview changes
- 🧩 Easy to extend file rules

## 📦 Supported File Categories

- Images
  `.jpg .jpeg .png .webp .tiff .tif .psd .raw .avif .svg .gif`
- Videos
  `.mp4 .mkv .avi .webm .mov .flv .wmv`
- Documents
  `.pdf .doc .xls .ppt .docx .xlsx .pptx .csv .odt .odp .ods .txt`
- Archives
  `.zip .rar .7z .tar`
- Other
  All unknown file types will be moved to: `other/`

## 🚀 Installation

1. Clone the repository

```bash
git clone https://github.com/MEMUZ/cli-sorter.git
cd cli-sorter
```

2. Build the binary

```bash
 go build -o sorter.exe
```

## 🖥 Usage

Basic usage

```bash
sorter.exe <directory>
```

Example:

```bash
sorter.exe ~/Downloads
```

Result:

```
Downloads/
  images/
  videos/
  documents/
  archives/
  other/
```

## 🔍 Dry Run Mode

Preview what will happen without moving files.

```bash
sorter.exe <directory> --dry-run
```

Example:

```bash
sorter.exe ~/Downloads --dry-run
```

Output example:

```bash
[DRY] image.png -> images
[DRY] movie.mp4 -> videos
[DRY] file.pdf -> documents
```
