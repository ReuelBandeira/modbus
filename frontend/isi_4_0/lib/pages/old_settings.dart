import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/pages/old_save_connection_mqtt.dart';
import 'package:isi_4_0/pages/save_configuration.dart';
import 'package:isi_4_0/providers/cards.dart';
import 'package:isi_4_0/providers/lang.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/language.dart';
import 'package:isi_4_0/widgets/custom_card_horizontal_home.dart';
import 'package:provider/provider.dart';

class Configuration extends StatefulWidget {
  const Configuration({super.key});

  @override
  State<Configuration> createState() => _ConfigurationState();
}

class _ConfigurationState extends State<Configuration> {
  AuthService authService = AuthService.instance;

  late List<String> dataCard = List.filled(4, "");

  var buttonAction = 0;
  final double heightCard = 140;
  final double widthCard = 300;
  final pageController = PageController();

  @override
  void initState() {
    super.initState();
    update();
  }

  void update() {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    (context).select((CardProvider card) => card.updateAll());
    return NotificationListener<SizeChangedLayoutNotification>(
      onNotification: (notification) {
        return true;
      },
      child: SizeChangedLayoutNotifier(
        child: Consumer<LanguageProvider>(
          builder: (context, language, child) => LayoutBuilder(
            builder: (context, constraints) {
              // Warning: breaking language changes to old screens
              // List<String> cardTitle =
              //     Language().getListFromCardTitleJson(language.currentLanguage);
              List<String> cardTitle = [];
              dataCard[0] = (context)
                  .select((CardProvider card) => card.configuredDevices);
              dataCard[1] =
                  (context).select((CardProvider card) => card.usedMemory);
              dataCard[2] =
                  (context).select((CardProvider card) => card.diskUsed);
              dataCard[3] = (context)
                  .select((CardProvider card) => card.connectedDevices);
              return Scaffold(
                backgroundColor: CustomColors.whiteColorLow,
                body: OverflowBox(
                  minHeight: 0,
                  minWidth: 0,
                  maxHeight: constraints.maxHeight,
                  maxWidth: constraints.maxWidth,
                  child: Column(
                    children: [
                      SizedBox(
                        height: heightCard,
                        width: constraints.maxWidth,
                        child: ListView.builder(
                          scrollDirection: Axis.horizontal,
                          itemCount: cardTitle.length,
                          itemBuilder: (context, index) =>
                              Consumer<CardProvider>(
                            builder: (context, card, child) =>
                                CustomCardHorizontalHome(
                              heightCard: heightCard,
                              widthCard: widthCard,
                              backgroundColor: CustomColors.whiteColorHigh,
                              cardName: cardTitle[index],
                              cardData: dataCard[index],
                              cardIcon: Icons.cached_outlined,
                              onTap: () {
                                switch (index) {
                                  case 0:
                                    card.updateConfiguredDevices();
                                    break;
                                  case 1:
                                    card.updateUsedMemory();
                                    break;
                                  case 2:
                                    card.updateDiskUsed();
                                    break;
                                  case 3:
                                    card.updateConnectedDevices();
                                    break;
                                  default:
                                }
                              },
                              shadowColor: CustomColors.greenColorLight,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(
                        height: 16,
                      ),
                      Container(
                        margin: const EdgeInsets.only(left: 2),
                        height: constraints.maxHeight - heightCard - 16,
                        width: constraints.maxWidth,
                        child: const SingleChildScrollView(
                          scrollDirection: Axis.vertical,
                          child: Wrap(
                            children: [
                              SaveConfiguration(),
                              SaveConnectionMqtt()
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}
