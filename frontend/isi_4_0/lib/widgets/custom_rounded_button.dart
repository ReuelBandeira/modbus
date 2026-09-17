import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class CustomRoundedButton extends StatelessWidget {
  final String textName;
  final double height;
  final double width;
  final double fontSize;
  final bool isSelected;
  final Color textColorActived;
  final Color textColorInactive;
  final Color backgroundColorActived;
  final Color backgroundColorInactive;
  final Color splashColor;
  final double borderRadiusValue;
  final Function()? onTap;
  final double? marginTop;
  const CustomRoundedButton(
      {Key? key,
      required this.textName,
      required this.height,
      required this.width,
      required this.fontSize,
      required this.isSelected,
      required this.textColorActived,
      required this.textColorInactive,
      required this.splashColor,
      required this.backgroundColorActived,
      required this.backgroundColorInactive,
      required this.borderRadiusValue,
      required this.onTap,
      this.marginTop})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      constraints: BoxConstraints(
        minWidth: width,
      ),
      margin: marginTop != null ?
        EdgeInsets.only(top: marginTop!, left: 1.0, right: 1.0, bottom: 1.0) :
        const EdgeInsets.all(1),
      decoration: BoxDecoration(
        color: Colors.transparent,
        borderRadius: BorderRadius.all(Radius.circular(borderRadiusValue)),
        boxShadow: const [
          BoxShadow(
            color: Color.fromARGB(127, 158, 158, 158),
            blurRadius: 4,
          )
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.all(Radius.circular(borderRadiusValue)),
        child: Material(
          color: isSelected
              ? backgroundColorActived
              : backgroundColorInactive, // button color
          child: InkWell(
            splashColor: splashColor,
            onTap: onTap, // button pressed
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                FittedBox(
                  fit: BoxFit.fill,
                  child: Container(
                    alignment: Alignment.center,
                    height: height,
                    padding: const EdgeInsets.only(left: 24, right: 24),
                    child: Text(
                      textName,
                      style: GoogleFonts.openSans(
                          fontSize: fontSize,
                          color: isSelected
                              ? textColorActived
                              : textColorInactive),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ), // text
              ],
            ),
          ),
        ),
      ),
    );
  }
}
