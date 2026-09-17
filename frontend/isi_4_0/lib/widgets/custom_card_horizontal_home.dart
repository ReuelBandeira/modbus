import 'package:flutter/material.dart';

import '../colors/color.dart';

class CustomCardHorizontalHome extends StatefulWidget {
  final VoidCallback onTap;
  final double widthCard;
  final double heightCard;
  final Color backgroundColor;
  final Color shadowColor;
  final String cardName;
  final IconData cardIcon;
  final String cardData;
  const CustomCardHorizontalHome(
      {super.key,
      required this.widthCard,
      required this.heightCard,
      required this.shadowColor,
      required this.cardName,
      required this.cardData,
      required this.cardIcon,
      required this.backgroundColor,
      required this.onTap});

  @override
  State<CustomCardHorizontalHome> createState() =>
      _CustomCardHorizontalHomeState();
}

class _CustomCardHorizontalHomeState extends State<CustomCardHorizontalHome> {
  @override
  Widget build(BuildContext context) {
    late Color iconColor = CustomColors.whiteColorHigh;
    return Padding(
      padding: const EdgeInsets.all(2.0),
      child: SizedBox(
        height: widget.heightCard,
        width: widget.widthCard,
        child: Card(
          elevation: 1,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(2.0),
          ),
          child: Stack(
            children: [
              Align(
                alignment: Alignment.topLeft,
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Text(
                    widget.cardName,
                    style: const TextStyle(color: Colors.black),
                  ),
                ),
              ),
              Align(
                alignment: Alignment.centerRight,
                child: Padding(
                  padding: const EdgeInsets.all(12.0),
                  child: InkWell(
                    onTap: widget.onTap,
                    borderRadius: BorderRadius.circular(45),
                    highlightColor: CustomColors.secondaryColorApp,
                    hoverColor: Colors.grey.withOpacity(0.3),
                    splashColor: Colors.transparent,
                    child: Container(
                      margin: const EdgeInsets.all(2),
                      height: 40,
                      width: 40,
                      decoration: BoxDecoration(
                        color: iconColor,
                        borderRadius: BorderRadius.circular(90),
                      ),
                      child: Icon(
                        widget.cardIcon,
                        color: CustomColors.secondaryColorApp,
                      ),
                    ),
                  ),
                ),
              ),
              Align(
                alignment: Alignment.bottomLeft,
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Text(
                    widget.cardData,
                    style: const TextStyle(
                      color: Colors.black,
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
