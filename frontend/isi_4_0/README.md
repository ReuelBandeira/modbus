# isi_4_0

A new Flutter project.

## Getting Started

This project is a starting point for a Flutter application.

A few resources to get you started if this is your first Flutter project:

- [Lab: Write your first Flutter app](https://docs.flutter.dev/get-started/codelab)
- [Cookbook: Useful Flutter samples](https://docs.flutter.dev/cookbook)

For help getting started with Flutter development, view the
[online documentation](https://docs.flutter.dev/), which offers tutorials,
samples, guidance on mobile development, and a full API reference.

## To install dependency use the follow commands:

```
 flutter pub add google_fonts
  flutter pub add url_launcher
  flutter pub add image_picker
  flutter pub add responsive_framework
  flutter pub add http
  flutter pub add fluttertoast
  flutter pub add font_awesome_flutter
```

### Run Modbus Emulator:

Open a Command Prompt on Windows, or Terminal on Linux, or Terminal on Linux Open folder **cts_finep_isi_40/bin-linux** or on Windows Open folder **cts_finep_isi_40/bin-windows** and type the following command :

```
./DeviceEmulator
  Select Option 14
    Select Option 2
```

### Run Modbus Isi 4.0 Emulator:

Open another Command Prompt on Windows, or Terminal on Linux, or Terminal on Linux Open folder **cts_finep_isi_40/bin-linux** or on Windows Open folder **cts_finep_isi_40/bin-windows** and type the following command :

```
./DeviceEmulator
  Select Option 14
    Select Option 3
```



### Run Frontend Web Browser:

```
flutter run -d chrome --web-renderer html
```

It will open a Browser and you will be asked to Login then going to other options

## To build application run the follow commands:

```
flutter build web --web-renderer html
```

### Errors

Error: The method 'File.create' has fewer named arguments than those of overridden method 'File.create'. Future<File> create({bool recursive = false});

Solve:
flutter clean
flutter pub upgrade
