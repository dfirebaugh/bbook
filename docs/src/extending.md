
# Extending

## Adding Additional JS to the Configuration

You can enhance the functionality and styling of your site by adding custom JavaScript (JS) and Cascading Style Sheets (CSS) files. To do this, use the `additional-js` and `additional-css` fields in your `book.toml` configuration file.

## Steps to Include Additional JS 

1. **Locate the Configuration File:**
   Ensure you're editing the `book.toml` file in your project's root directory.

2. **Specify JS Files:**
   Under the `[output.html]` section, add an `additional-js` field. List the paths to your JavaScript files relative to the source directory. For example:

```toml
[output.html]
additional-js = ["js/custom-script.js", "js/another-script.js"]
```


## Adding Additional CSS to the Configuration

To include custom CSS files in your project, you can use the `additional-css` field in the `book.toml` configuration file. Follow these steps to add CSS files:

1. **Locate the Configuration File:**
   Open the `book.toml` file in your project's root directory.

2. **Specify CSS Files:**
   Under the `[output.html]` section, add an `additional-css` field. List the paths to your CSS files relative to the source directory. For example:

```toml
[output.html]
additional-css = ["css/custom-style.css", "css/another-style.css"]
```
