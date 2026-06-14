#import <AppKit/AppKit.h>
#include <stdlib.h>
#include <stdint.h>

void goMenuCallback(int menuID);
void goOpenFileCallback(const char *path);

@interface MDAppDelegate : NSObject <NSApplicationDelegate>
- (void)setupMenuWithLang:(NSString *)lang;
@property (nonatomic, strong) NSString *currentLang;
@property (nonatomic, assign) NSUInteger currentLangIndex;
@property (nonatomic, strong) NSString *pendingFile;
@property (nonatomic, strong) NSPanel *aboutPanel;
@property (nonatomic, strong) NSTextField *aboutNameLabel;
@property (nonatomic, strong) NSTextField *aboutVersionLabel;
@property (nonatomic, strong) NSTextField *aboutCopyrightLabel;
@property (nonatomic, strong) NSImageView *aboutIconView;
@end

@implementation MDAppDelegate

- (void)applicationWillFinishLaunching:(NSNotification *)notification {
    // This helps ensure we catch openFile: events during startup
}

- (BOOL)application:(NSApplication *)sender openFile:(NSString *)filename {
    if (filename.length == 0) {
        return NO;
    }
    self.pendingFile = filename;
    goOpenFileCallback([filename UTF8String]);
    return YES;
}

/// macOS delivers file opens via URLs in addition to openFile:; implement both.
- (void)application:(NSApplication *)application openURLs:(NSArray<NSURL *> *)urls API_AVAILABLE(macos(10.13)) {
    for (NSURL *url in urls) {
        if (![url isFileURL]) {
            continue;
        }
        NSString *path = [url path];
        if (path.length == 0) {
            continue;
        }
        self.pendingFile = path;
        goOpenFileCallback([path UTF8String]);
    }
}

- (void)application:(NSApplication *)sender openFiles:(NSArray<NSString *> *)filenames {
    for (NSString *file in filenames) {
        if (file.length == 0) {
            continue;
        }
        self.pendingFile = file;
        goOpenFileCallback([file UTF8String]);
    }
}

- (void)setupMenuWithLang:(NSString *)lang {
    self.currentLang = lang.length > 0 ? lang : @"zhTW";

    // i18n dictionary: key -> array [en, zhTW, zhCN, ja, ko]
    NSDictionary *tr = @{
        @"menuFile":      @[@"File", @"檔案", @"文件", @"ファイル", @"파일"],
        @"menuEdit":      @[@"Edit", @"編輯", @"编辑", @"編集", @"편집"],
        @"editCopy":      @[@"Copy", @"拷貝", @"拷贝", @"コピー", @"복사"],
        @"editPaste":     @[@"Paste", @"貼上", @"粘贴", @"貼り付け", @"붙여넣기"],
        @"editSelectAll": @[@"Select All", @"全選", @"全选", @"すべてを選択", @"모두 선택"],
        @"editWritingTools": @[@"Writing Tools…", @"書寫工具…", @"书写工具…", @"ライティングツール…", @"작성 도구…"],
        @"menuView":      @[@"View", @"顯示", @"显示", @"表示", @"보기"],
        @"menuExport":    @[@"Export", @"匯出", @"导出", @"エクスポート", @"내보내기"],
        @"menuHelp":      @[@"Help", @"說明", @"幫助", @"ヘルプ", @"도움말"],
        @"appAbout":      @[@"About md-viewer", @"關於 md-viewer", @"关于 md-viewer", @"md-viewer について", @"md-viewer 정보"],
        @"appPref":       @[@"Preferences...", @"偏好設定...", @"偏好设置...", @"環境設定...", @"설정..."],
        @"appQuit":       @[@"Quit md-viewer", @"結束 md-viewer", @"结束 md-viewer", @"md-viewer を終了", @"md-viewer 종료"],
        @"fileOpen":      @[@"Open...", @"開啟檔案...", @"打开文件...", @"ファイルを開く...", @"파일 열기..."],
        @"fileReload":    @[@"Reload", @"重新載入", @"重新载入", @"再読み込み", @"새로고침"],
        @"recentOpen":    @[@"Open Recent", @"最近開啟", @"最近打开", @"最近開いたファイル", @"최근에 연 파일"],
        @"recentEmpty":   @[@"No Recent Files", @"無最近檔案", @"无最近文件", @"最近ファイルなし", @"최근 파일 없음"],
        @"recentTooltip": @[@"File not found, click to remove", @"檔案不存在，點擊移除", @"文件不存在，点击移除", @"ファイルが見つかりません。クリックして削除", @"파일을 찾을 수 없습니다. 클릭하여 제거"],
        @"fileExportHTML":@[@"Export as HTML...", @"匯出為 HTML...", @"导出为 HTML...", @"HTML としてエクスポート...", @"HTML로 내보내기..."],
        @"fileExportPDF": @[@"Export as PDF...",  @"匯出為 PDF...",  @"导出为 PDF...",  @"PDF としてエクスポート...",  @"PDF로 내보내기..."],
        @"viewIn":        @[@"Zoom In", @"放大", @"放大", @"拡大", @"확대"],
        @"viewOut":       @[@"Zoom Out", @"縮小", @"缩小", @"縮小", @"축소"],
        @"viewReset":     @[@"Actual Size", @"實際大小", @"实际大小", @"實際のサイズ", @"실제 크기"],
        @"viewFull":      @[@"Toggle Full Screen", @"切換全螢幕", @"切换全屏", @"フルスクリーン切替", @"전체 화면 전환"],
        @"viewFocus":     @[@"Focus Mode", @"專注模式", @"专注模式", @"集中モード", @"집중 모드"],
        @"viewFind":      @[@"Find…", @"尋找…", @"查找…", @"検索…", @"찾기…"],
        @"viewTranslate":  @[@"Translate...", @"翻譯文件...", @"翻译文档...", @"翻訳...", @"번역..."],
        @"helpAbout":     @[@"About md-viewer", @"關於 md-viewer", @"关于 md-viewer", @"md-viewer について", @"md-viewer 정보"],
    };
    NSArray *langs = @[@"en", @"zhTW", @"zhCN", @"ja", @"ko"];
    NSUInteger li = [langs indexOfObject:self.currentLang];
    if (li == NSNotFound) li = 0;

    NSString *(^t)(NSString *) = ^(NSString *key) {
        NSArray *arr = tr[key];
        return arr ? arr[li] : key;
    };
    
    self.currentLangIndex = li;

    // Rebuild main menu
    NSMenu *mainMenu = [[NSMenu alloc] init];

    // App menu
    NSMenuItem *appItem = [[NSMenuItem alloc] init];
    [appItem setTitle:@"md-viewer"];
    NSMenu *appMenu = [[NSMenu alloc] init];
    NSMenuItem *aboutAppItem = [[NSMenuItem alloc] initWithTitle:t(@"appAbout") action:@selector(mdShowAboutPanel) keyEquivalent:@""];
    [aboutAppItem setTarget:self];
    [appMenu addItem:aboutAppItem];
    [appMenu addItem:[NSMenuItem separatorItem]];
    // ⌘, toggles settings panel (standard macOS shortcut)
    NSMenuItem *prefShortcut = [[NSMenuItem alloc] initWithTitle:t(@"appPref") action:@selector(menuPreferences) keyEquivalent:@","];
    [appMenu addItem:prefShortcut];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:t(@"appQuit") action:@selector(menuQuit) keyEquivalent:@"q"];
    [appItem setSubmenu:appMenu];
    [mainMenu addItem:appItem];

    // File menu
    NSMenuItem *fileItem = [[NSMenuItem alloc] init];
    [fileItem setTitle:t(@"menuFile")];
    NSMenu *fileMenu = [[NSMenu alloc] initWithTitle:t(@"menuFile")];
    [fileMenu addItemWithTitle:t(@"fileOpen") action:@selector(menuOpen) keyEquivalent:@"o"];
    [fileMenu addItemWithTitle:t(@"fileReload") action:@selector(menuReload) keyEquivalent:@"r"];
    [fileMenu addItem:[NSMenuItem separatorItem]];
    
    // Recent files submenu
    NSMenuItem *recentItem = [[NSMenuItem alloc] init];
    [recentItem setTitle:t(@"recentOpen")];
    NSMenu *recentMenu = [[NSMenu alloc] init];
    recentItem.submenu = recentMenu;
    [fileMenu addItem:recentItem];
    
    // Export submenu
    NSMenuItem *exportItem = [[NSMenuItem alloc] init];
    [exportItem setTitle:t(@"menuExport")];
    NSMenu *exportMenu = [[NSMenu alloc] init];
    [exportMenu addItemWithTitle:t(@"fileExportHTML") action:@selector(menuExportHTML) keyEquivalent:@""];
    [exportMenu addItemWithTitle:t(@"fileExportPDF") action:@selector(menuExportPDF) keyEquivalent:@""];
    [exportItem setSubmenu:exportMenu];
    [fileMenu addItem:exportItem];
    [fileItem setSubmenu:fileMenu];
    [mainMenu addItem:fileItem];

    // Edit menu: minimal items so ⌘C / ⌘V / ⌘A reach WKWebView (nil target → responder chain)
    NSMenuItem *editItem = [[NSMenuItem alloc] init];
    [editItem setTitle:t(@"menuEdit")];
    NSMenu *editMenu = [[NSMenu alloc] initWithTitle:t(@"menuEdit")];
    [editMenu addItemWithTitle:t(@"editCopy") action:@selector(copy:) keyEquivalent:@"c"];
    [editMenu addItemWithTitle:t(@"editPaste") action:@selector(paste:) keyEquivalent:@"v"];
    [editMenu addItemWithTitle:t(@"editSelectAll") action:@selector(selectAll:) keyEquivalent:@"a"];
    if (@available(macOS 15.2, *)) {
        [editMenu addItem:[NSMenuItem separatorItem]];
        NSMenuItem *wtItem = [[NSMenuItem alloc] initWithTitle:t(@"editWritingTools")
                                                        action:NSSelectorFromString(@"showWritingTools:")
                                                 keyEquivalent:@""];
        [editMenu addItem:wtItem];
    }
    [editItem setSubmenu:editMenu];
    [mainMenu addItem:editItem];

    // View menu
    NSMenuItem *viewItem = [[NSMenuItem alloc] init];
    [viewItem setTitle:t(@"menuView")];
    NSMenu *viewMenu = [[NSMenu alloc] initWithTitle:t(@"menuView")];
    [viewMenu addItemWithTitle:t(@"viewIn") action:@selector(menuZoomIn) keyEquivalent:@"="];
    [viewMenu addItemWithTitle:t(@"viewOut") action:@selector(menuZoomOut) keyEquivalent:@"-"];
    [viewMenu addItemWithTitle:t(@"viewReset") action:@selector(menuZoomReset) keyEquivalent:@"0"];
    [viewMenu addItemWithTitle:t(@"viewTOC") action:@selector(menuToggleTOC) keyEquivalent:@"t"];
    [viewMenu addItem:[NSMenuItem separatorItem]];
    NSMenuItem *translateItem = [[NSMenuItem alloc] initWithTitle:t(@"viewTranslate") action:@selector(menuTranslate) keyEquivalent:@"t"];
    translateItem.keyEquivalentModifierMask = NSEventModifierFlagCommand | NSEventModifierFlagShift;
    [viewMenu addItem:translateItem];
    [viewMenu addItem:[NSMenuItem separatorItem]];
    NSMenuItem *focusItem = [[NSMenuItem alloc] initWithTitle:t(@"viewFocus") action:@selector(menuFocusMode) keyEquivalent:@"m"];
    focusItem.keyEquivalentModifierMask = NSEventModifierFlagCommand | NSEventModifierFlagShift;
    [viewMenu addItem:focusItem];
    NSMenuItem *findItem = [[NSMenuItem alloc] initWithTitle:t(@"viewFind") action:@selector(menuFind) keyEquivalent:@"f"];
    findItem.keyEquivalentModifierMask = NSEventModifierFlagCommand | NSEventModifierFlagShift;
    [viewMenu addItem:findItem];
    [viewMenu addItem:[NSMenuItem separatorItem]];
    [viewMenu addItemWithTitle:t(@"viewFull") action:@selector(menuFullscreen) keyEquivalent:@"f"];
    [viewItem setSubmenu:viewMenu];
    [mainMenu addItem:viewItem];

    // Help menu
    NSMenuItem *helpItem = [[NSMenuItem alloc] init];
    [helpItem setTitle:t(@"menuHelp")];
    NSMenu *helpMenu = [[NSMenu alloc] initWithTitle:t(@"menuHelp")];
    NSMenuItem *aboutHelpItem = [[NSMenuItem alloc] initWithTitle:t(@"helpAbout") action:@selector(mdShowAboutPanel) keyEquivalent:@""];
    [aboutHelpItem setTarget:self];
    [helpMenu addItem:aboutHelpItem];
    [helpItem setSubmenu:helpMenu];
    [mainMenu addItem:helpItem];

    [NSApp setMainMenu:mainMenu];
}

/// Custom About window with larger type (system About panel text is too small on some displays).
- (void)mdShowAboutPanel {
    [self mdEnsureAboutPanel];
    [self mdRefreshAboutPanelStrings];
    [self.aboutPanel center];
    [NSApp activateIgnoringOtherApps:YES];
    [self.aboutPanel makeKeyAndOrderFront:nil];
}

- (void)mdEnsureAboutPanel {
    if (self.aboutPanel) {
        return;
    }

    NSRect rect = NSMakeRect(0, 0, 480, 320);
    NSPanel *panel = [[NSPanel alloc] initWithContentRect:rect
                                                  styleMask:(NSWindowStyleMaskTitled | NSWindowStyleMaskClosable)
                                                    backing:NSBackingStoreBuffered
                                                      defer:NO];
    panel.title = @"About md-viewer";
    panel.releasedWhenClosed = NO;
    panel.level = NSFloatingWindowLevel;

    NSImageView *iconView = [[NSImageView alloc] initWithFrame:NSZeroRect];
    iconView.image = [NSApp applicationIconImage];
    iconView.imageScaling = NSImageScaleProportionallyDown;
    iconView.translatesAutoresizingMaskIntoConstraints = NO;

    NSTextField *nameLabel = [NSTextField labelWithString:@"md-viewer"];
    nameLabel.font = [NSFont boldSystemFontOfSize:22];
    nameLabel.alignment = NSTextAlignmentCenter;
    nameLabel.translatesAutoresizingMaskIntoConstraints = NO;

    NSTextField *verLabel = [NSTextField labelWithString:@"0.0.0"];
    verLabel.font = [NSFont monospacedDigitSystemFontOfSize:17 weight:NSFontWeightMedium];
    verLabel.textColor = [NSColor labelColor];
    verLabel.alignment = NSTextAlignmentCenter;
    verLabel.translatesAutoresizingMaskIntoConstraints = NO;

    NSTextField *copyLabel = [NSTextField wrappingLabelWithString:@""];
    copyLabel.font = [NSFont systemFontOfSize:15];
    copyLabel.textColor = [NSColor secondaryLabelColor];
    copyLabel.alignment = NSTextAlignmentCenter;
    copyLabel.maximumNumberOfLines = 0;
    copyLabel.preferredMaxLayoutWidth = 400;
    copyLabel.translatesAutoresizingMaskIntoConstraints = NO;

    NSStackView *stack = [NSStackView stackViewWithViews:@[iconView, nameLabel, verLabel, copyLabel]];
    stack.orientation = NSUserInterfaceLayoutOrientationVertical;
    stack.alignment = NSLayoutAttributeCenterX;
    stack.spacing = 16;
    stack.edgeInsets = NSEdgeInsetsMake(28, 28, 28, 28);
    stack.translatesAutoresizingMaskIntoConstraints = NO;

    NSView *content = panel.contentView;
    [content addSubview:stack];
    [NSLayoutConstraint activateConstraints:@[
        [stack.leadingAnchor constraintEqualToAnchor:content.leadingAnchor],
        [stack.trailingAnchor constraintEqualToAnchor:content.trailingAnchor],
        [stack.topAnchor constraintEqualToAnchor:content.topAnchor],
        [stack.bottomAnchor constraintEqualToAnchor:content.bottomAnchor],
        [iconView.widthAnchor constraintEqualToConstant:80],
        [iconView.heightAnchor constraintEqualToConstant:80],
    ]];

    self.aboutPanel = panel;
    self.aboutIconView = iconView;
    self.aboutNameLabel = nameLabel;
    self.aboutVersionLabel = verLabel;
    self.aboutCopyrightLabel = copyLabel;
}

- (void)mdRefreshAboutPanelStrings {
    NSBundle *bundle = [NSBundle mainBundle];
    NSDictionary *info = [bundle infoDictionary] ?: @{};

    NSString *name = info[@"CFBundleName"];
    if (name.length == 0) {
        name = @"md-viewer";
    }
    NSString *shortVersion = info[@"CFBundleShortVersionString"];
    NSString *build = info[@"CFBundleVersion"];
    NSString *copyright = info[@"NSHumanReadableCopyright"];

    if (shortVersion.length == 0 && build.length == 0) {
        shortVersion = @"0.0.0-dev";
        build = @"";
        if (copyright.length == 0) {
            copyright = @"請使用 ./build.sh 打包後以 md-viewer.app 開啟，以顯示完整版本資訊。";
        }
    }

    NSMutableString *verLine = [NSMutableString string];
    if (shortVersion.length > 0) {
        [verLine appendString:shortVersion];
    }
    if (build.length > 0) {
        if (verLine.length > 0) {
            [verLine appendFormat:@"  ·  %@", build];
        } else {
            [verLine appendString:build];
        }
    }
    if (verLine.length == 0) {
        [verLine appendString:@"—"];
    }

    self.aboutNameLabel.stringValue = name;
    self.aboutVersionLabel.stringValue = verLine;
    self.aboutCopyrightLabel.stringValue = copyright ?: @"";
    self.aboutCopyrightLabel.hidden = (copyright.length == 0);
}

- (void)menuPreferences { goMenuCallback(2); }
- (void)menuOpen        { goMenuCallback(3); }
- (void)menuReload     { goMenuCallback(4); }
- (void)menuOpenRecent:(NSMenuItem *)sender {
    NSString *path = sender.representedObject;
    if (path) {
        goOpenFileCallback([path UTF8String]);
    }
}
void goRemoveRecentFileCallback(const char *path);

- (void)menuRemoveRecent:(NSMenuItem *)sender {
    NSString *path = sender.representedObject;
    if (path) {
        goRemoveRecentFileCallback([path UTF8String]);
    }
}
- (void)menuQuit        { [NSApp terminate:nil]; }
- (void)menuZoomIn      { goMenuCallback(6); }
- (void)menuZoomOut     { goMenuCallback(7); }
- (void)menuZoomReset   { goMenuCallback(8); }
- (void)menuToggleTOC   { goMenuCallback(15); }
- (void)menuTranslate   { NSLog(@"[menu.m] menuTranslate called"); goMenuCallback(16); }
- (void)menuExportHTML  { goMenuCallback(12); }
- (void)menuExportPDF   { goMenuCallback(13); }
- (void)menuFocusMode  { goMenuCallback(14); }
- (void)menuFind       { goMenuCallback(17); }
- (void)menuFullscreen  {
    NSWindow *window = [NSApp keyWindow];
    if (window && [window respondsToSelector:@selector(toggleFullScreen:)]) {
        [window toggleFullScreen:nil];
    }
}

@end

static MDAppDelegate *_sharedDelegate = nil;
static NSString *currentMenuLang = @"";

/// Call before webview.New so application:openFile: has a delegate (avoids
/// "cannot open files in the … format" while the file still opens via argv).
void EnsureNSAppDelegateInstalled(void) {
    [NSApplication sharedApplication];
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        if (_sharedDelegate == nil) {
            _sharedDelegate = [[MDAppDelegate alloc] init];
        }
        [NSApp setDelegate:_sharedDelegate];
    });
}

/// After Go registers openFileCallback, replay any file Apple delivered earlier.
void FlushPendingOpenFile(void) {
    MDAppDelegate *d = _sharedDelegate;
    if (!d || d.pendingFile.length == 0) {
        return;
    }
    NSString *path = [d.pendingFile copy];
    d.pendingFile = nil;
    goOpenFileCallback([path UTF8String]);
}

void UpdateMenuLanguageTitles(const char *lang) {
    currentMenuLang = [NSString stringWithUTF8String:lang];
    dispatch_async(dispatch_get_main_queue(), ^{
        if (_sharedDelegate) {
            [(MDAppDelegate *)_sharedDelegate setupMenuWithLang:currentMenuLang];
        }
    });
}

void SetupMainMenu(void) {
    EnsureNSAppDelegateInstalled();
    static dispatch_once_t menuOnce;
    dispatch_once(&menuOnce, ^{
        NSString *lang = currentMenuLang.length > 0 ? currentMenuLang : @"zhTW";
        [(MDAppDelegate *)_sharedDelegate setupMenuWithLang:lang];
    });
}

static NSDictionary *RecentTranslations(void) {
    return @{
        @"recentOpen":    @[@"Open Recent", @"最近開啟", @"最近打开", @"最近開いたファイル", @"최근에 연 파일"],
        @"recentEmpty":   @[@"No Recent Files", @"無最近檔案", @"无最近文件", @"最近ファイルなし", @"최근 파일 없음"],
        @"recentTooltip": @[@"File not found, click to remove", @"檔案不存在，點擊移除", @"文件不存在，点击移除", @"ファイルが見つかりません。クリックして削除", @"파일을 찾을 수 없습니다. 클릭하여 제거"],
    };
}

static NSArray *MenuItemTitles(void) {
    return @[@"File", @"檔案", @"文件", @"ファイル", @"파일"];
}

NSString *TranslateRecent(NSString *key, NSUInteger langIndex) {
    NSArray *arr = RecentTranslations()[key];
    return arr ? arr[langIndex] : key;
}

void UpdateRecentFilesMenu(const char **files, int count) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSMenu *mainMenu = [NSApp mainMenu];
        if (!mainMenu) return;
        
        NSUInteger langIndex = 1;
        if (_sharedDelegate) {
            langIndex = [(MDAppDelegate *)_sharedDelegate currentLangIndex];
        }
        
        NSArray *items = mainMenu.itemArray;
        NSMenuItem *fileItem = nil;
        for (NSMenuItem *item in items) {
            for (NSString *title in MenuItemTitles()) {
                if ([item.title isEqualToString:title]) {
                    fileItem = item;
                    break;
                }
            }
            if (fileItem) break;
        }
        if (!fileItem) return;
        NSMenu *fileMenu = fileItem.submenu;
        if (!fileMenu) return;
        
        NSMenuItem *recentItem = nil;
        NSArray *recentTitles = RecentTranslations()[@"recentOpen"];
        for (int i = 0; i < fileMenu.numberOfItems; i++) {
            NSMenuItem *item = [fileMenu itemAtIndex:i];
            for (NSString *title in recentTitles) {
                if ([item.title isEqualToString:title]) {
                    recentItem = item;
                    break;
                }
            }
            if (recentItem) break;
        }
        
        if (!recentItem) return;
        
        recentItem.title = TranslateRecent(@"recentOpen", langIndex);
        
        NSMenu *recentMenu = recentItem.submenu;
        [recentMenu removeAllItems];
        
        for (int i = 0; i < count; i++) {
            NSString *filePath = [NSString stringWithUTF8String:files[i]];
            if (!filePath) continue;
            
            BOOL fileExists = [[NSFileManager defaultManager] fileExistsAtPath:filePath];
            NSString *fileName = [filePath lastPathComponent];
            
            NSMenuItem *item = [[NSMenuItem alloc] initWithTitle:fileName 
                                                           action:fileExists ? @selector(menuOpenRecent:) : @selector(menuRemoveRecent:) 
                                                    keyEquivalent:@""];
            item.representedObject = [filePath copy];
            
            if (!fileExists) {
                NSMutableAttributedString *attrTitle = [[NSMutableAttributedString alloc] initWithString:fileName];
                [attrTitle addAttribute:NSForegroundColorAttributeName 
                                   value:[NSColor disabledControlTextColor] 
                                   range:NSMakeRange(0, fileName.length)];
                item.attributedTitle = attrTitle;
                item.toolTip = TranslateRecent(@"recentTooltip", langIndex);
            }
            
            [recentMenu addItem:item];
        }
        
        if (count == 0) {
            NSMenuItem *emptyItem = [[NSMenuItem alloc] initWithTitle:TranslateRecent(@"recentEmpty", langIndex)
                                                               action:nil 
                                                        keyEquivalent:@""];
            emptyItem.enabled = NO;
            [recentMenu addItem:emptyItem];
        }
    });
}

// Set window frame (position + size)
void SetWindowFrame(void *windowPtr, int x, int y, int width, int height) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSWindow *window = (__bridge NSWindow *)windowPtr;
        if (window) {
            NSRect frame = NSMakeRect(x, y, width, height);
            [window setFrame:frame display:YES animate:NO];
        }
    });
}

// Get current window size
void GetWindowSize(void *windowPtr, int *width, int *height) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (window) {
        NSRect frame = [window frame];
        *width = (int)frame.size.width;
        *height = (int)frame.size.height;
    } else {
        *width = 0;
        *height = 0;
    }
}

// Get current window position
void GetWindowPosition(void *windowPtr, int *x, int *y) {
    NSWindow *window = (__bridge NSWindow *)windowPtr;
    if (window) {
        NSRect frame = [window frame];
        *x = (int)frame.origin.x;
        *y = (int)frame.origin.y;
    } else {
        *x = 0;
        *y = 0;
    }
}
