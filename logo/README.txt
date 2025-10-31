Favicon.pub 文件包说明
===================

本压缩包包含以下文件：

PNG 格式：
- favicon-16x16.png (16x16像素，用于浏览器标签页)
- favicon-32x32.png (32x32像素，用于任务栏)
- favicon-48x48.png (48x48像素，用于桌面快捷方式)
- android-chrome-192x192.png (192x192像素,Android PWA图标)
- android-chrome-512x512.png (512x512像素,Android PWA启动图标)

Apple 图标：
- apple-touch-icon.png (180x180像素,iOS默认图标)
- apple-touch-icon-152x152.png (152x152像素,iPad图标)
- apple-touch-icon-167x167.png (167x167像素,iPad Pro图标)
- apple-touch-icon-180x180.png (180x180像素,iPhone图标)

其它格式：
- favicon.ico (传统ICO格式，包含 16x16, 32x32, 48x48 多尺寸。请注意，浏览器或操作系统可能会根据显示环境选择最合适的尺寸，或因缓存原因显示较小尺寸。您可以使用专业的ICO查看器（如 Axialis IconWorkshop 或在线ICO检查工具）来验证文件中是否确实嵌入了所有尺寸。)
- favicon.svg (现代SVG格式，目前为通用占位符。请注意：在浏览器端将位图（如 PNG, JPG）高质量地转换为矢量图（SVG）是一个非常复杂的任务，超出了当前客户端 JavaScript 的能力。因此，此文件目前仅为占位符。建议您手动创建或使用专业的矢量图形软件（如 Adobe Illustrator, Inkscape）将您的 Logo 转换为 SVG 格式，以获得最佳效果。)
- safari-pinned-tab.svg (Safari固定标签页图标，目前为通用占位符。同上，建议手动创建矢量版本。)
- site.webmanifest (PWA清单文件)

使用方法：
将所有文件放置在网站根目录，然后在HTML的<head>标签中添加相应的链接标签。
