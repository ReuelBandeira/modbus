import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';

class NavigationService {
  static final NavigationService instance = NavigationService._internal();
  factory NavigationService() {
    return instance;
  }
  NavigationService._internal();

  final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

  Future<dynamic> navigateTo(String routeName) {
    return navigatorKey.currentState!.pushNamed(routeName);
  }

  Future<dynamic> logOut() {
    AuthService().removeLogin();
    AuthService().refreshPage('/');
    return navigatorKey.currentState!.pushNamed('/');
  }
}
