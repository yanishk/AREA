import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../custom_widget.dart";
import "../globalData.dart";

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(
              Icons.arrow_back,
              size: 32,
              color: color,
            ),
          ),
        ),
        Center(
          child: AutoSizeText(
            title,
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: color,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}

class SelectActionOrReactionTopBar extends StatelessWidget {
  const SelectActionOrReactionTopBar({
    super.key,
    required this.title,
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(
              Icons.arrow_back,
              size: 32,
              color: color,
            ),
          ),
        ),
        Center (
          child: AutoSizeText(
            title,
            style: TextStyle(
              fontSize: 25,
              fontWeight: FontWeight.bold,
              color: color,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}

class ActionOrReactionDescription extends StatelessWidget {
  const ActionOrReactionDescription({
    super.key,
    required this.colorTheme,
    required this.title,
    required this.serviceIcon,
    required this.serviceName,
    required this.actionDescription,
  });

  final Color colorTheme;
  final String title;
  final String serviceIcon;
  final String serviceName;
  final String actionDescription;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: colorTheme,
        border: Border.all(
          color: colorTheme,
          width: 2.0,
        ),
        borderRadius: BorderRadius.circular(4.0),
      ),
      child: InkWell(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectActionOrReactionTopBar(title: title, color: Colors.white),
            Padding(padding: EdgeInsets.only(top: 20)),
            SafeWebImage(
              url:
                  "${Globaldata.domainName}$serviceIcon",
              width: 75,
              height: 75,
            ),
            Padding(padding: EdgeInsets.only(top: 5)),
            AutoSizeText(
              serviceName,
              style: TextStyle(
                fontSize: 25,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
              maxLines: 1,
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 15.0),
              child: AutoSizeText(
                actionDescription,
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
                textAlign: TextAlign.center,
                maxLines: 10,
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
          ],
        ),
      ),
    );
  }
}

class SelectDateTimeWidget extends StatelessWidget {
  const SelectDateTimeWidget({
    super.key,
    required this.text,
    this.colorTheme = Colors.black,
    required this.onPressed,
  });

  final String text;
  final Color colorTheme;
  final Function onPressed;

  @override
  Widget build(BuildContext context) {
    return CustomButton(
      backgroundColor: colorTheme,
      text: AutoSizeText(
        text,
        style: TextStyle(
            fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
      ),
      onPressed: () {onPressed;},
      size: Size(double.infinity, 60),
      padding: EdgeInsets.symmetric(horizontal: 15.0),
    );
  }
}
