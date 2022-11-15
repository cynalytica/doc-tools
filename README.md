# Doc Tools

## How to Doc

File structure of docs should look like this:

```text
- {whatever file name for docs section of project, ie "docs"}
    - files <- all md files
        - file1.md
        - file2.md
        - fileDir
            - file3.md
        - etc...
    - media <- whatever supporting files, images, etc
        - image1.jpg
        - imageDir
            - image2.jpg
            - image3.jpg
    - manifest.json <- generated file for TOC
```

### How to TOC

The `manifest.json` file is what maps out resources for the UI. The UI will build a TOC based on the `manifest.json` file.

The `manifest.json` file is updated automatically - any manual changes will not be persisted.

All links to pages in UI should have an absolute route. All links internal to this set up docs should be relative including images.
- `"See [Actions](../actions) for more info"`
- `"![image title](../media/image1.jpg 'alt text')"`

Images won't work in dev because the proxy isn't redirecting traffic, you can temporarily change them to absolute path (including host) at appropriate port to see how they look in dev in the UI.

### Env Vars

TBD. For now, see `flags/Flags.go`.

### Serving Docs

All files in the directories should be served statically. For example if serving a markdown file from `/api/docs`, files would be served like this:
- `/api/docs/files/file1.md`
- `/api/docs/media/imageDir/image2.jpg`
- `/api/docs/manifest.json`

The UI will request each file and rewrite/prettify the path.

## How to PDF

### Install:
- pandoc
- xelatex
  - also several plugins required
  - full list TBD


### Powershell command:

```
pandoc -s $(Get-ChildItem -Path .\ -Filter *.md) --toc --pdf-engine=xelatex --from markdown+escaped_line_breaks+backtick_code_blocks+pipe_tables+multiline_tables+fenced_code_attributes --title "CyRenQL Documentation" --syntax-definition "./cyrenql.xml"  -o "./../doc.pdf" --metadata "title=CyRenQL Documentation" --metadata "subtitle=Get the most out the Cynalytica AnalytICS Engine" --template "./cytemplate.tex"
```

### Pandoc

Helpful Pandoc docs: https://pandoc.org/MANUAL.html