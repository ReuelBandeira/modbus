import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/model/cards_model.dart';

class CustomCards extends StatelessWidget {
  const CustomCards({super.key, required this.info});
  final InfoCards info;
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 120,
      width: 200,
      margin: const EdgeInsets.only(right: 8.0),
      child: Card(
        elevation: 1,
        color: CustomColors.whiteColorLow,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(2.0),
        ),
        child: Padding(
          padding: const EdgeInsets.only(left: 8.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                  alignment: Alignment.bottomLeft,
                  padding: const EdgeInsets.only(bottom: 8.0),
                  child: Text(info.nameCard,
                      style: const TextStyle(
                          fontSize: 13, color: CustomColors.primaryColorApp))),
              Container(
                  alignment: Alignment.topLeft,
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text(info.dataCard,
                      style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                          color: CustomColors.primaryColorApp))),
            ],
          ),
        ),
      ),
    );
  }
}
