# PDFExtract
A Go Command Line App to get all the source URLS from an arxiv pdf.

# Demo
![alt text][logo]

[logo]: https://github.com/aaditkapoor/PDFExtract/blob/master/demo.gif

# Setup
- The application uses xpdf-tools so make sure to download the command line tools before getting started from here: http://www.xpdfreader.com/download.html
- The application requires that you have a pdfs/ directory in the app. The app also automatically creates the folder for you.

# Quickstart
- Clone the repo and run go build.
- To run: ./PDFExtract -url https://arxiv.org/pdf/2005.14187v1.pdf -download true or go run main.go ...

# Future Direction
- To add recommendation support.
- Ability to download code
- Web interface

# LICENSE
Copyright 2020 Aadit Kapoor

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
