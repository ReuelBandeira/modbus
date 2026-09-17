import 'package:flutter/material.dart';
import 'package:isi_4_0/views/home/home.dart';
// import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/viewmodel/home_view_model.dart';
import 'package:isi_4_0/views/login/recover.dart';
import 'package:isi_4_0/views/login/register.dart';
import 'package:isi_4_0/views/login/sign_in.dart';
import 'package:isi_4_0/viewmodel/login_view_model.dart';
import 'package:provider/provider.dart';

class Routes {
  static Route<dynamic>? generateRoute(RouteSettings routes) {
    switch (routes.name) {
      case '/':
        return MaterialPageRoute(
          builder: (_) => MultiProvider(
            providers: [
              ChangeNotifierProvider(create: (_) => LoginViewModel()),
              // ChangeNotifierProvider(create: (_) => LanguageProvider()),
            ],
            child: const SignIn(),
          ),
        );
      case '/recover':
        return MaterialPageRoute(
          builder: (_) => MultiProvider(
            providers: [
              ChangeNotifierProvider(create: (_) => LoginViewModel()),
              // ChangeNotifierProvider(create: (_) => LanguageProvider()),
            ],
            child: const Recover(),
          ),
        );
      case '/register':
        return MaterialPageRoute(
          builder: (_) => MultiProvider(
            providers: [
              ChangeNotifierProvider(create: (_) => LoginViewModel()),
              // ChangeNotifierProvider(create: (_) => LanguageProvider()),
            ],
            child: const Register(),
          ),
        );
      case '/home':
        return MaterialPageRoute(
          builder: (_) => MultiProvider(
            providers: [
              ChangeNotifierProvider(create: (_) => HomeViewModel()),
              // ChangeNotifierProvider(create: (_) => LanguageProvider()),
            ],
            child: const Home(),
          ),
        );
      default:
    }
    return null;
  }
}
