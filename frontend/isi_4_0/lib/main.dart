import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/cards.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/providers/profile.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/routes/routes.dart';
import 'package:isi_4_0/utils/app_details.dart';
import 'package:isi_4_0/utils/navigator_service.dart';
import 'package:provider/provider.dart';
import 'package:toastification/toastification.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setPreferredOrientations([DeviceOrientation.landscapeLeft]).then(
    (_) {
      runApp(
        MultiProvider(
          providers: [
            ChangeNotifierProvider(create: (_) => LanguageProvider()),
            ChangeNotifierProvider(create: (_) => CardProvider()),
            ChangeNotifierProvider(create: (_) => ProfileProvider()),
          ],
          child: const Isi(),
        ),
      );
    },
  );
}

class Isi extends StatelessWidget {
  const Isi({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    AuthService service = AuthService.instance;
    return FutureBuilder(
      future: service.getInitialPage(),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.done) {
          return MaterialApp(
            debugShowCheckedModeBanner: false,
            title: AppDetails.name.text,
            theme: ThemeData(
                appBarTheme:
                    const AppBarTheme(color: CustomColors.primaryColorApp),
                colorScheme: ThemeData().colorScheme.copyWith(
                    primary: CustomColors.primaryColorApp,
                    secondary: CustomColors.secondaryColorApp,
                    background: CustomColors.whiteColorHigh),
                tooltipTheme: TooltipThemeData(
                    decoration: BoxDecoration(
                        color: CustomColors.primaryColorApp.withOpacity(0.75),
                        borderRadius: BorderRadius.circular(90))),
                textTheme: GoogleFonts.robotoTextTheme(
                  Theme.of(context).textTheme,
                )),
            onGenerateRoute: Routes.generateRoute,
            navigatorKey: NavigationService().navigatorKey,
            initialRoute: snapshot.data.toString(),
            builder: (context, child) {
              return ToastificationConfigProvider(
                config: ToastificationConfig(
                  marginBuilder: (geo) {
                      return const EdgeInsets.fromLTRB(0, 60, 0, 0);
                    },
                  animationDuration: const Duration(milliseconds: 300),
                ),
                child: child!,
              );
            },
          );
        } else {
          // Make loading page here
          return Container();
        }
      },
    );
  }
}
