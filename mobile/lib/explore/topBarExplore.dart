import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "explorePageData.dart";

class TopBarExplore extends StatelessWidget {
  TopBarExplore({
    super.key,
    required this.onPressed,
  });

  final Function onPressed;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {
            SelectedCardList.cardType = SelectedCardListEnum.All;
            onPressed();
          },
          child: TextTopBar(text: "All", cardType: SelectedCardListEnum.All),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.cardType = SelectedCardListEnum.Workflows;
            onPressed();
          },
          child: TextTopBar(text: "Workflows", cardType: SelectedCardListEnum.Workflows,),
        ),
        TextButton(
          onPressed: () {
            SelectedCardList.cardType = SelectedCardListEnum.Services;
            onPressed();
          },
          child: TextTopBar(text: "Services", cardType: SelectedCardListEnum.Services,),
        ),
      ],
    );
  }
}

class TextTopBar extends StatelessWidget {
  const TextTopBar({
    super.key,
    required this.text,
    this.cardType = SelectedCardListEnum.All,
  });

  final String text;
  final SelectedCardListEnum cardType;

  @override
  Widget build(BuildContext context) {
    if (cardType != SelectedCardList.cardType) {
      return AutoSizeText(
        text,
        style: TextStyle(
            fontSize: 20, fontWeight: FontWeight.w600, color: Colors.black),
        maxLines: 1,
      );
    }
    return AutoSizeText(
        text,
        style: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.w900,
          color: Colors.black,
          decoration: TextDecoration.underline,
          decorationThickness: 2,
        ),
        maxLines: 1,
      );
  }
}
